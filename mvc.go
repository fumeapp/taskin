package taskin

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (r *Runners) Init() tea.Cmd {
	var cmds []tea.Cmd
	for i := range *r {
		cmds = append(cmds, (*r)[i].Spinner.Tick)
	}
	return tea.Batch(cmds...)
}

func (r *Runners) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case spinner.TickMsg:
		allDone := true
		for i := range *r {

			if (*r)[i].State == Running || (*r)[i].State == NotStarted {
				newSpinner, cmd := (*r)[i].Spinner.Update(msg)
				(*r)[i].Spinner = newSpinner
				cmds = append(cmds, cmd)
			}

			if (*r)[i].State != Completed && (*r)[i].State != Failed {
				allDone = false
			}
		}

		if allDone {
			return r, tea.Quit
		}

		return r, tea.Batch(cmds...)
	}

	return r, nil
}

func (r *Runners) View() string {
	var view string
	for _, runner := range *r {
		status := ""
		switch runner.State {
		case NotStarted:
			status = lipgloss.NewStyle().Foreground(runner.Config.Colors.Pending).Render("•") + " " + runner.Task.Title // Gray bullet
		case Running:
			if runner.Task.ShowProgress.Total != 0 {
				percent := float64(runner.Task.ShowProgress.Current) / float64(runner.Task.ShowProgress.Total)
				status = runner.Spinner.View() + runner.Task.Title + " " + runner.Task.Bar.ViewAs(percent)
			} else {
				status = runner.Spinner.View() + " " + runner.Task.Title
			}
		case Completed:
			status = lipgloss.NewStyle().Foreground(runner.Config.Colors.Success).Render("✔") + " " + runner.Task.Title // Green checkmark
		case Failed:
			status = lipgloss.NewStyle().Foreground(runner.Config.Colors.Failure).Render("✘") + " " + runner.Task.Title // Red 'x'
		}
		view += lipgloss.NewStyle().Render(status) + "\n"
	}
	return view
}
