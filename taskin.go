package taskin

import (
	"dario.cat/mergo"
	"fmt"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func NewRunner(task Task, cfg Config) Runner {
	s := spinner.New(spinner.WithSpinner(cfg.Spinner))           // Initialize with a spinner model
	s.Style = lipgloss.NewStyle().Foreground(cfg.Colors.Spinner) // Styling spinner
	return Runner{Task: task, State: NotStarted, Spinner: s, Config: cfg}
}

func (task *Task) Progress(current, total int) {
	task.ShowProgress = TaskProgress{Current: current, Total: total}
	if !task.Bar.IsAnimating() {
		task.Bar = progress.New(task.Config.ProgressOption)
	}
	if total != 0 { // Check if TaskProgress is set
		percent := float64(current) / float64(total)
		task.Bar.SetPercent(percent)
	}
}

func (r *Runners) Run() error {
	p := tea.NewProgram(r)
	_, err := p.Run()
	return err
}

func New(tasks Tasks, cfg Config) Runners {
	// merge cfg with Defaults
	_ = mergo.Merge(&cfg, Defaults)
	var runners Runners
	for _, task := range tasks {
		// Use NewRunner to ensure runners are initialized with spinners correctly
		task.Config = cfg
		runners = append(runners, NewRunner(task, cfg))
	}
	go func() {
		for i := range runners {
			runners[i].State = Running
			err := runners[i].Task.Task(&runners[i].Task)
			if err != nil {
				runners[i].Task.Title = fmt.Sprintf("%s - Error: %s", runners[i].Task.Title, err.Error())
				runners[i].State = Failed
				continue
			}
			runners[i].State = Completed
		}
	}()
	return runners
}
