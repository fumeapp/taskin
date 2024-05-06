package taskin

import (
	"dario.cat/mergo"
	"fmt"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os"
)

var program *tea.Program

func NewRunner(task Task, cfg Config) Runner {

	var spinr *spinner.Model

	if os.Getenv("CI") == "" {
		spinnerModel := spinner.New(spinner.WithSpinner(cfg.Spinner))           // Initialize with a spinner model
		spinnerModel.Style = lipgloss.NewStyle().Foreground(cfg.Colors.Spinner) // Styling spinner
		spinr = &spinnerModel
	}

	children := make(Runners, len(task.Tasks))
	for i, childTask := range task.Tasks {
		childTask.Config = cfg
		children[i] = NewRunner(childTask, cfg)
	}
	if task.Task == nil {
		task.Task = func(t *Task) error {
			return nil
		}
	}
	return Runner{Task: task, State: NotStarted, Spinner: spinr, Config: cfg, Children: children}
}

func (task *Task) Progress(current, total int) {
	task.ShowProgress = TaskProgress{Current: current, Total: total}
	if !task.Bar.IsAnimating() {
		task.Bar = progress.New(task.Config.ProgressOptions...)
	}
	if total != 0 { // Check if TaskProgress is set
		percent := float64(current) / float64(total)
		task.Bar.SetPercent(percent)
	}
}

func (r *Runners) Run() error {
	program = tea.NewProgram(r, tea.WithInput(nil))
	_, err := program.Run()
	return err
}

func New(tasks Tasks, cfg Config) Runners {
	_ = mergo.Merge(&cfg, Defaults)
	var runners Runners
	for _, task := range tasks {
		task.Config = cfg
		runners = append(runners, NewRunner(task, cfg))
	}

	go func() {
		for i := range runners {

			for _, runner := range runners[:i] {
				if runner.State == Failed {
					program.Send(spinner.TickMsg{})
					return
				}
			}

			runners[i].State = Running
			err := runners[i].Task.Task(&runners[i].Task)
			if err != nil {
				runners[i].Task.Title = fmt.Sprintf("%s - Error: %s", runners[i].Task.Title, err.Error())
				runners[i].State = Failed
				program.Send(spinner.TickMsg{})
				continue
			}

			// Run child tasks
			for j := range runners[i].Children {
				runners[i].Children[j].State = Running
				err := runners[i].Children[j].Task.Task(&runners[i].Children[j].Task)
				if err != nil {
					runners[i].Children[j].Task.Title = fmt.Sprintf("%s - Error: %s", runners[i].Children[j].Task.Title, err.Error())
					runners[i].Children[j].State = Failed
					runners[i].State = Failed // Mark parent task as Failed
					program.Send(spinner.TickMsg{})
					break
				}
				runners[i].Children[j].State = Completed
			}

			// Check if all child tasks are completed
			allChildrenCompleted := true
			for _, child := range runners[i].Children {
				if child.State != Completed {
					allChildrenCompleted = false
					break
				}
			}

			// If all child tasks are completed, mark the parent task as completed
			if allChildrenCompleted && runners[i].State != Failed {
				runners[i].State = Completed
				program.Send(spinner.TickMsg{})
			}
		}
	}()
	return runners
}
