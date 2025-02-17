package taskin

import (
	"dario.cat/mergo"
	"fmt"
	"io"
	"os"
	"regexp"

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

		if task.ShowProgress.Total != 0 {
			task.Bar = progress.New(cfg.ProgressOptions...)
		}
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

	var out io.Writer = os.Stdout
	if IsCI() {
		out = &ansiEscapeCodeFilter{writer: out}
	}

	program = tea.NewProgram(m, tea.WithInput(nil), tea.WithOutput(out))
	_, err := program.Run()
	if err != nil {
		return fmt.Errorf("program run error: %w", err)
	}

	if m.Shutdown && m.ShutdownError != nil {
		return fmt.Errorf("shutdown error: %w", m.ShutdownError)
	}

	return nil
}

func New(tasks Tasks, cfg Config) Runners {
	_ = mergo.Merge(&cfg, Defaults)
	var runners Runners
	for _, task := range tasks {
		task.Config = cfg
		runners = append(runners, NewRunner(task, cfg))
	}

	// Helper function to run a task and its children recursively
	var runTaskAndChildren func(runner *Runner) error
	runTaskAndChildren = func(runner *Runner) error {
		runner.State = Running

		// Run the task itself first if it has a function
		if runner.Task.Task != nil {
			err := runner.Task.Task(&runner.Task)
			if err != nil {
				runner.Task.Title = fmt.Sprintf("%s - %s", runner.Task.Title, err.Error())
				runner.State = Failed
				return err
			}
		}

		// Run all children recursively
		for i := range runner.Children {
			err := runTaskAndChildren(&runner.Children[i])
			if err != nil {
				runner.State = Failed
				return err
			}
		}

		runner.State = Completed
		if program != nil {
			program.Send(spinner.TickMsg{})
		}
		return nil
	}

	go func() {
		for i := range runners {
			// Check for previous failures
			for _, prev := range runners[:i] {
				if prev.State == Failed && prev.Config.Options.ExitOnFailure {
					return
				}
			}

			err := runTaskAndChildren(&runners[i])
			if err != nil && program != nil {
				program.Send(TerminateWithError{Error: err})
			}
		}
	}()

	return runners
}

func IsCI() bool {
	return os.Getenv("CI") != "" ||
		os.Getenv("CONTINUOUS_INTEGRATION") != "" ||
		os.Getenv("BUILD_NUMBER") != "" ||
		os.Getenv("GITHUB_ACTIONS") != ""
}
