package taskin

import (
	"fmt"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/octoper/go-ray"
)

func NewRunner(task Task, cfg Config) Runner {
	s := spinner.New(spinner.WithSpinner(cfg.Spinner))           // Initialize with a spinner model
	s.Style = lipgloss.NewStyle().Foreground(cfg.Colors.Spinner) // Styling spinner
	return Runner{Task: task, State: NotStarted, Spinner: s, Config: cfg}
}

func (task *Task) Progress(current, total int) {
	task.ShowProgress = TaskProgress{Current: current, Total: total}
	if !task.Bar.IsAnimating() {
		task.Bar = progress.New(progress.WithDefaultGradient())
	}
	if total != 0 { // Check if TaskProgress is set
		percent := float64(current) / float64(total)
		task.Bar.SetPercent(percent)
		ray.Ray(percent)
	}
}

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

func (r *Runners) Run() error {
	p := tea.NewProgram(r)
	_, err := p.Run()
	return err
}

func New(tasks Tasks, cfg Config) Runners {
	var runners Runners
	for _, task := range tasks {
		// Use NewRunner to ensure runners are initialized with spinners correctly
		runners = append(runners, NewRunner(task, cfg))
	}
	// simulate tasks running after the program starts
	go func() {
		for i := range runners {
			runners[i].State = Running
			err := runners[i].Task.Task(&runners[i].Task)
			if err != nil {
				runners[i].Task.Title = fmt.Sprintf("%s - Error: %s", runners[i].Task.Title, err.Error())
				runners[i].State = Failed
				continue
			}
			runners[i].State = Completed // Example of updating state, replace with real task logic
		}
	}()
	return runners
}
