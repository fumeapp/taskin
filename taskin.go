package taskin

import (
	"fmt"
	"io"
	"os"
	"regexp"

	"dario.cat/mergo"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var program *tea.Program

func NewRunner(task Task, cfg Config) Runner {

	var spinr *spinner.Model

	if !IsCI() {
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
	if IsCI() {
		return
	}
	if !task.Bar.IsAnimating() {
		task.Bar = progress.New(task.Config.ProgressOptions...)
	}
	if total != 0 { // Check if TaskProgress is set
		percent := float64(current) / float64(total)
		task.Bar.SetPercent(percent)
	}
}

type ansiEscapeCodeFilter struct {
	writer io.Writer
}

func (f *ansiEscapeCodeFilter) Write(p []byte) (n int, err error) {
	// Corrected regular expression to match ANSI escape codes
	re := regexp.MustCompile(` *\x1b\[[0-?]*[ -/]*[@-~]`)
	// Remove the escape codes from the input
	p = re.ReplaceAll(p, []byte{})
	// Write the filtered input to the original writer
	return f.writer.Write(p)
}

func (r *Runners) Run() error {
	m := &Model{Runners: *r, Shutdown: false, ShutdownError: nil}
	if IsCI() {
		program = tea.NewProgram(m, tea.WithInput(nil), tea.WithOutput(&ansiEscapeCodeFilter{writer: os.Stdout}))
	} else {
		program = tea.NewProgram(m, tea.WithInput(nil))
	}
	_, err := program.Run()
	if err != nil {
		program.Send(TerminateWithError{Error: err})
	}
	if m.Shutdown && m.ShutdownError != nil {
		return m.ShutdownError
	}
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
				if runner.State == Failed && runner.Config.Options.ExitOnFailure {
					return
				}
			}

			runners[i].State = Running
			err := runners[i].Task.Task(&runners[i].Task)
			if err != nil {
				runners[i].Task.Title = fmt.Sprintf("%s - %s", runners[i].Task.Title, err.Error())
				runners[i].State = Failed
				if program != nil {
					program.Send(TerminateWithError{Error: err})
				}
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
					if program != nil {
						program.Send(TerminateWithError{Error: err})
					}
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
				if program != nil {
					program.Send(spinner.TickMsg{})
				}
			}
		}
	}()
	return runners
}

func IsCI() bool {
	return os.Getenv("CI") != ""
}
