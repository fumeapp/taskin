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
	return tea.Batch(cmds...)
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

	case spinner.TickMsg:
		allDone := true
		for i := range m.Runners {

			if (m.Runners)[i].State == Running || (m.Runners)[i].State == NotStarted {
				if !IsCI() {
					newSpinner, cmd := (m.Runners)[i].Spinner.Update(msg)
					(m.Runners)[i].Spinner = &newSpinner
					cmds = append(cmds, cmd)
				}
			}

			for j := range (m.Runners)[i].Children {
				if (m.Runners)[i].Children[j].State == Running || (m.Runners)[i].Children[j].State == NotStarted {
					if !IsCI() {
						newSpinner, cmd := (m.Runners)[i].Children[j].Spinner.Update(msg)
						(m.Runners)[i].Children[j].Spinner = &newSpinner
						cmds = append(cmds, cmd)
					}
				}
			}

			if (m.Runners)[i].State == Failed {
				return m, tea.Quit
			}

			if (m.Runners)[i].State != Completed && (m.Runners)[i].State != Failed {
				allDone = false
			}
		}

		if allDone {
			return m, tea.Quit
		}

		return m, tea.Batch(cmds...)
	}

	return m, nil
}

func (m *Model) SetShutdown(err error) {
	m.Shutdown = true
	m.ShutdownError = err
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

	// Handle CI mode
	if IsCI() {
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

// In mvc.go, update the renderTask function:

func renderTask(runner Runner, indent string) string {
	var view string
	status := ""

	switch runner.State {
	case NotStarted:
		status = Color(runner.Config.Colors.Pending, runner.Config.Chars.NotStarted) + " " + runner.Task.Title
	case Running:
		if len(runner.Children) > 0 {
			status = Color(runner.Config.Colors.ParentStarted, runner.Config.Chars.ParentStarted) + " " + runner.Task.Title
		} else {
			// Handle both progress bars and spinners consistently for all levels
			if runner.Task.ShowProgress.Total != 0 {
				percent := float64(runner.Task.ShowProgress.Current) / float64(runner.Task.ShowProgress.Total)
				if runner.Spinner != nil {
					status = runner.Spinner.View() + " " + runner.Task.Title + " " + runner.Task.Bar.ViewAs(percent)
				}
			} else if runner.Spinner != nil {
				status = runner.Spinner.View() + " " + runner.Task.Title
			}
		}
	case Completed:
		status = Color(runner.Config.Colors.Success, runner.Config.Chars.Success) + " " + runner.Task.Title
	case Failed:
		status = Color(runner.Config.Colors.Failure, runner.Config.Chars.Failure) + " " + runner.Task.Title
	}

	if IsCI() {
		view += indent + status + "\n"
	} else {
		view += indent + lipgloss.NewStyle().Render(status) + "\n"
	}

	// Recursively render children if parent is running
	if len(runner.Children) > 0 && (runner.State == Running || IsCI()) {
		for _, child := range runner.Children {
			view += renderTask(child, indent+"  ")
		}
	}

	return view
}
