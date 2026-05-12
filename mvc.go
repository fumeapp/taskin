package taskin

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m *Model) Init() tea.Cmd {

	var cmds []tea.Cmd
	for i := range m.Runners {
		if (m.Runners)[i].Spinner != nil {
			cmds = append(cmds, (m.Runners)[i].Spinner.Tick)
		}
		for j := range (m.Runners)[i].Children {
			if (m.Runners)[i].Children[j].Spinner != nil {
				cmds = append(cmds, (m.Runners)[i].Children[j].Spinner.Tick)
			}
		}
	}
	if m.taskMessages != nil {
		cmds = append(cmds, runTasksCmd(cloneRunners(m.Runners), m.taskMessages), waitForTaskMessage(m.taskMessages))
	}
	return tea.Batch(cmds...)
}

func waitForTaskMessage(messages <-chan tea.Msg) tea.Cmd {
	return func() tea.Msg {
		return <-messages
	}
}

func (m *Model) waitForTaskMessage() tea.Cmd {
	if m.taskMessages == nil {
		return nil
	}
	return waitForTaskMessage(m.taskMessages)
}

func (m *Model) SetShutdown(err error) {
	m.Shutdown = true
	m.ShutdownError = err
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	if m.Shutdown && m.ShutdownError != nil {
		return m, tea.Quit
	}

	switch msg := msg.(type) {
	case TerminateWithError:
		m.SetShutdown(msg.Error)
		return m, tea.Quit
	case taskExecutionFinishedMsg:
		return m, tea.Quit
	case taskStartedMsg:
		if runner := m.runnerAtPath(msg.Path); runner != nil {
			runner.State = Running
		}
		return m, m.waitForTaskMessage()
	case taskUpdatedMsg:
		if runner := m.runnerAtPath(msg.Path); runner != nil {
			runner.Task = msg.Task
		}
		return m, m.waitForTaskMessage()
	case taskCompletedMsg:
		if runner := m.runnerAtPath(msg.Path); runner != nil {
			runner.Task = msg.Task
			runner.State = Completed
		}
		allDone, anyFailed := m.checkTasksState()
		if allDone && !anyFailed {
			return m, tea.Quit
		}
		return m, m.waitForTaskMessage()
	case taskFailedMsg:
		if runner := m.runnerAtPath(msg.Path); runner != nil {
			runner.Task = msg.Task
			runner.State = Failed
		}
		return m, m.waitForTaskMessage()

	case spinner.TickMsg:
		// Helper function to update spinners recursively
		var updateSpinners func(runner *Runner) []tea.Cmd
		updateSpinners = func(runner *Runner) []tea.Cmd {
			var spinnerCmds []tea.Cmd

			if runner.State == Running || runner.State == NotStarted {
				if !IsCI() && !runner.Config.DisableUI && runner.Spinner != nil {
					newSpinner, cmd := runner.Spinner.Update(msg)
					runner.Spinner = &newSpinner
					spinnerCmds = append(spinnerCmds, cmd)
				}
			}

			// Recursively update all children's spinners
			for i := range runner.Children {
				spinnerCmds = append(spinnerCmds, updateSpinners(&runner.Children[i])...)
			}

			return spinnerCmds
		}

		allDone := true
		for i := range m.Runners {
			cmds = append(cmds, updateSpinners(&m.Runners[i])...)

			if m.Runners[i].State != Completed && m.Runners[i].State != Failed {
				allDone = false
			}
		}

		_, anyFailed := m.checkTasksState()
		if allDone && !anyFailed {
			return m, tea.Quit
		}

		return m, tea.Batch(cmds...)
	}

	return m, nil
}

func (m *Model) runnerAtPath(path []int) *Runner {
	if len(path) == 0 || path[0] < 0 || path[0] >= len(m.Runners) {
		return nil
	}
	runner := &m.Runners[path[0]]
	for _, index := range path[1:] {
		if index < 0 || index >= len(runner.Children) {
			return nil
		}
		runner = &runner.Children[index]
	}
	return runner
}

func (m *Model) checkTasksState() (allDone, anyFailed bool) {
	allDone = true
	for _, runner := range m.Runners {
		if runner.State != Completed && runner.State != Failed {
			allDone = false
		}
		if runner.State == Failed {
			anyFailed = true
		}
	}
	return
}

func (m *Model) View() string {
	// Check for hidden views first
	for _, runner := range m.Runners {
		if runner.Task.HideView {
			return ""
		}
	}

	// Handle CI mode or UI disabled mode
	if IsCI() || (len(m.Runners) > 0 && m.Runners[0].Config.DisableUI) {
		allDone, anyFailed := m.checkTasksState()
		if !allDone && !anyFailed {
			return ""
		}
	}

	var view string
	for _, runner := range m.Runners {
		view += renderTask(runner, "")
	}
	return view
}

func renderTask(runner Runner, indent string) string {
	var view string
	status := ""

	switch runner.State {
	case NotStarted:
		status = Color(runner.Config.Colors.Pending, runner.Config.Chars.NotStarted) + " " + runner.Task.Title
	case Running:
		if len(runner.Children) > 0 {
			status = Color(runner.Config.Colors.ParentStarted, runner.Config.Chars.ParentStarted) + " " + runner.Task.Title
		} else if runner.Task.ShowProgress.Total != 0 {
			percent := float64(runner.Task.ShowProgress.Current) / float64(runner.Task.ShowProgress.Total)
			status = runner.Spinner.View() + runner.Task.Title + " " + runner.Task.Bar.ViewAs(percent)
		} else if runner.Spinner != nil {
			status = runner.Spinner.View() + runner.Task.Title
		}
	case Completed:
		status = Color(runner.Config.Colors.Success, runner.Config.Chars.Success) + " " + runner.Task.Title
	case Failed:
		status = Color(runner.Config.Colors.Failure, runner.Config.Chars.Failure) + " " + runner.Task.Title
	}

	if IsCI() || runner.Config.DisableUI {
		view = indent + status + "\n"
	} else {
		view = indent + lipgloss.NewStyle().Render(status) + "\n"
	}

	// Recursively render children
	if len(runner.Children) > 0 && (runner.State == Running || IsCI() || runner.Config.DisableUI) {
		for _, child := range runner.Children {
			view += renderTask(child, indent+"  ")
		}
	}

	return view
}
