package taskin

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
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

	case spinner.TickMsg:
		// Helper function to update spinners recursively
		var updateSpinners func(runner *Runner) []tea.Cmd
		updateSpinners = func(runner *Runner) []tea.Cmd {
			var spinnerCmds []tea.Cmd

			if runner.State == Running || runner.State == NotStarted {
				if !IsCI() && runner.Spinner != nil {
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

			if m.Runners[i].State == Failed {
				return m, tea.Quit
			}

			if m.Runners[i].State != Completed && m.Runners[i].State != Failed {
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
