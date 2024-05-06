package taskin

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os"
)

func (r *Runners) Init() tea.Cmd {
	var cmds []tea.Cmd
	for i := range *r {
		if (*r)[i].Spinner != nil {
			cmds = append(cmds, (*r)[i].Spinner.Tick)
		}
		for j := range (*r)[i].Children {
			if (*r)[i].Children[j].Spinner != nil {
				cmds = append(cmds, (*r)[i].Children[j].Spinner.Tick)
			}
		}
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
				if os.Getenv("CI") == "" {
					newSpinner, cmd := (*r)[i].Spinner.Update(msg)
					(*r)[i].Spinner = &newSpinner
					cmds = append(cmds, cmd)
				}
			}

			for j := range (*r)[i].Children {
				if (*r)[i].Children[j].State == Running || (*r)[i].Children[j].State == NotStarted {
					if os.Getenv("CI") == "" {
						newSpinner, cmd := (*r)[i].Children[j].Spinner.Update(msg)
						(*r)[i].Children[j].Spinner = &newSpinner
						cmds = append(cmds, cmd)
					}
				}
			}

			if (*r)[i].State == Failed {
				return r, tea.Quit
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
