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

// In mvc.go, update the renderTask function:

func renderChildren(runner Runner, indent string) string {
	var view string
	for _, child := range runner.Children {
		status := ""

		// Only process children when parent is running or in CI mode
		if runner.State == Running || IsCI() {
			switch child.State {
			case NotStarted:
				status = Color(child.Config.Colors.Pending, child.Config.Chars.NotStarted) + " " + child.Task.Title
			case Running:
				if child.Task.ShowProgress.Total != 0 {
					percent := float64(child.Task.ShowProgress.Current) / float64(child.Task.ShowProgress.Total)
					status = child.Spinner.View() + " " + child.Task.Title + " " + child.Task.Bar.ViewAs(percent)
				} else {
					status = child.Spinner.View() + " " + child.Task.Title
				}
			case Completed:
				status = Color(child.Config.Colors.Success, child.Config.Chars.Success) + " " + child.Task.Title
			case Failed:
				status = Color(child.Config.Colors.Failure, child.Config.Chars.Failure) + " " + child.Task.Title
			}

			if IsCI() {
				view += indent + "  " + status + "\n"
			} else {
				view += indent + "  " + lipgloss.NewStyle().Render(status) + "\n"
			}

			// Recursively render any nested children
			if len(child.Children) > 0 {
				view += renderChildren(child, indent+"  ")
			}
		}
	}
	return view
}

// Update the View method to use renderChildren
func (m *Model) View() string {
	// ... existing view code for hiding and CI checks ...

	var view string
	for _, runner := range m.Runners {
		status := ""
		switch runner.State {
		case NotStarted:
			status = Color(runner.Config.Colors.Pending, runner.Config.Chars.NotStarted) + " " + runner.Task.Title
		case Running:
			if runner.Task.ShowProgress.Total != 0 {
				percent := float64(runner.Task.ShowProgress.Current) / float64(runner.Task.ShowProgress.Total)
				status = runner.Spinner.View() + " " + runner.Task.Title + " " + runner.Task.Bar.ViewAs(percent)
			} else {
				status = runner.Spinner.View() + " " + runner.Task.Title
			}
		case Completed:
			status = Color(runner.Config.Colors.Success, runner.Config.Chars.Success) + " " + runner.Task.Title
		case Failed:
			status = Color(runner.Config.Colors.Failure, runner.Config.Chars.Failure) + " " + runner.Task.Title
		}

		view += lipgloss.NewStyle().Render(status) + "\n"

		// Render children if they exist
		if len(runner.Children) > 0 {
			view += renderChildren(runner, "")
		}
	}
	return view
}
