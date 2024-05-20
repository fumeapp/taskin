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
	var view string

	// check if CI is set, if it is then don't return the view until all tasks are completed or one has failed
	if IsCI() {
		allDone, anyFailed := m.checkTasksState()
		if !allDone && !anyFailed {
			return ""
		}
	}

	for _, runner := range m.Runners {
		status := ""
		switch runner.State {
		case NotStarted:
			status = lipgloss.NewStyle().Foreground(runner.Config.Colors.Pending).Render(runner.Config.Chars.NotStarted) + " " + runner.Task.Title // Gray bullet
		case Running:
			if len(runner.Children) > 0 {
				status = lipgloss.NewStyle().Foreground(runner.Config.Colors.ParentStarted).Render(runner.Config.Chars.ParentStarted) + " " + runner.Task.Title
			} else {
				if runner.Task.ShowProgress.Total != 0 {
					percent := float64(runner.Task.ShowProgress.Current) / float64(runner.Task.ShowProgress.Total)
					if runner.Spinner != nil {
						status = runner.Spinner.View() + runner.Task.Title + " " + runner.Task.Bar.ViewAs(percent)
					} else {
						status = "⣟ " + runner.Task.Title + " " + runner.Task.Bar.ViewAs(percent)
					}
				} else {
					if runner.Spinner != nil {
						status = runner.Spinner.View() + runner.Task.Title
					} else {
						status = "⣟ " + runner.Task.Title
					}
				}
			}
		case Completed:
			status = lipgloss.NewStyle().Foreground(runner.Config.Colors.Success).Render(runner.Config.Chars.Success) + " " + runner.Task.Title // Green checkmark
		case Failed:
			status = lipgloss.NewStyle().Foreground(runner.Config.Colors.Failure).Render(runner.Config.Chars.Failure) + " " + runner.Task.Title // Red 'x'
		}
		view += lipgloss.NewStyle().Render(status) + "\n"

		for _, child := range runner.Children {
			status = ""
			switch child.State {
			case NotStarted:
				status = lipgloss.NewStyle().Foreground(child.Config.Colors.Pending).Render(runner.Config.Chars.NotStarted) + " " + child.Task.Title // Gray bullet
			case Running:
				if child.Task.ShowProgress.Total != 0 {
					percent := float64(child.Task.ShowProgress.Current) / float64(child.Task.ShowProgress.Total)
					if child.Spinner == nil {
						status = "⣟ " + child.Task.Title + " " + child.Task.Bar.ViewAs(percent)
					} else {
						status = child.Spinner.View() + child.Task.Title + " " + child.Task.Bar.ViewAs(percent)
					}
				} else {
					if child.Spinner == nil {
						status = "⣟ " + child.Task.Title
					} else {
						status = child.Spinner.View() + child.Task.Title
					}
				}
			case Completed:
				status = lipgloss.NewStyle().Foreground(child.Config.Colors.Success).Render("✔") + " " + child.Task.Title // Green checkmark
			case Failed:
				status = lipgloss.NewStyle().Foreground(child.Config.Colors.Failure).Render("✘") + " " + child.Task.Title // Red 'x'
			}
			view += "  " + lipgloss.NewStyle().Render(status) + "\n"
		}
	}
	return view
}
