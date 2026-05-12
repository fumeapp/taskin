package taskin

import (
	"dario.cat/mergo"
	"fmt"
	"io"
	"os"
	"regexp"
	"sync"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type taskUpdateCallback func(Task)

var taskUpdateCallbacks sync.Map

func NewRunner(task Task, cfg Config) Runner {

	var spinr *spinner.Model

	if !IsCI() && !cfg.DisableUI {
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
	task.applyProgress(TaskProgress{Current: current, Total: total})
	task.notifyUpdate()
}

func (task *Task) SetTitle(title string) {
	task.Title = title
	task.notifyUpdate()
}

func (task *Task) SetHideView(hide bool) {
	task.HideView = hide
	task.notifyUpdate()
}

func (task *Task) applyProgress(taskProgress TaskProgress) {
	task.ShowProgress = taskProgress
	if IsCI() || task.Config.DisableUI {
		return
	}
	if !task.Bar.IsAnimating() {
		task.Bar = progress.New(task.Config.ProgressOptions...)
	}
	if taskProgress.Total != 0 { // Check if TaskProgress is set
		percent := float64(taskProgress.Current) / float64(taskProgress.Total)
		task.Bar.SetPercent(percent)
	}
}

func (task *Task) notifyUpdate() {
	callback, ok := taskUpdateCallbacks.Load(task)
	if !ok {
		return
	}
	callback.(taskUpdateCallback)(snapshotTask(*task))
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
	m := &Model{Runners: cloneRunners(*r), Shutdown: false, ShutdownError: nil, taskMessages: make(chan tea.Msg, 64)}

	var out io.Writer = os.Stdout
	// Check if we need to disable UI features or are in CI mode
	if IsCI() || (len(*r) > 0 && (*r)[0].Config.DisableUI) {
		out = &ansiEscapeCodeFilter{writer: out}
	}

	program := tea.NewProgram(m, tea.WithInput(nil), tea.WithOutput(out))
	_, err := program.Run()
	*r = m.Runners
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

	return runners
}

func runTasksCmd(runners Runners, messages chan<- tea.Msg) tea.Cmd {
	return func() tea.Msg {
		go func() {
			if err := runRunners(runners, messages); err != nil {
				messages <- TerminateWithError{Error: err}
			}
		}()
		return nil
	}
}

func runRunners(runners Runners, messages chan<- tea.Msg) error {
	var firstErr error
	for i := range runners {
		err := runTaskAndChildren(&runners[i], pathWithIndex(nil, i), messages)
		if err == nil {
			continue
		}
		if firstErr == nil {
			firstErr = err
		}
		if runners[i].Config.Options.ExitOnFailure {
			return firstErr
		}
	}
	return firstErr
}

func runTaskAndChildren(runner *Runner, path []int, messages chan<- tea.Msg) error {
	runner.State = Running
	messages <- taskStartedMsg{Path: path}

	task := runner.Task
	callback := taskUpdateCallback(func(updated Task) {
		messages <- taskUpdatedMsg{Path: path, Task: updated}
	})
	taskUpdateCallbacks.Store(&task, callback)
	var err error
	if task.Task != nil {
		err = task.Task(&task)
	}
	taskUpdateCallbacks.Delete(&task)

	runner.Task = snapshotTask(task)
	if err != nil {
		runner.Task.Title = fmt.Sprintf("%s - %s", runner.Task.Title, err.Error())
		runner.State = Failed
		messages <- taskFailedMsg{Path: path, Task: runner.Task}
		return err
	}
	messages <- taskUpdatedMsg{Path: path, Task: runner.Task}

	for i := range runner.Children {
		if err := runTaskAndChildren(&runner.Children[i], pathWithIndex(path, i), messages); err != nil {
			runner.State = Failed
			messages <- taskFailedMsg{Path: path, Task: runner.Task}
			return err
		}
	}

	runner.State = Completed
	messages <- taskCompletedMsg{Path: path, Task: runner.Task}
	return nil
}

func cloneRunners(runners Runners) Runners {
	if runners == nil {
		return nil
	}
	cloned := make(Runners, len(runners))
	for i := range runners {
		cloned[i] = runners[i]
		cloned[i].Task = snapshotTask(runners[i].Task)
		cloned[i].Children = cloneRunners(runners[i].Children)
		if runners[i].Spinner != nil {
			spinnerCopy := *runners[i].Spinner
			cloned[i].Spinner = &spinnerCopy
		}
	}
	return cloned
}

func snapshotTask(task Task) Task {
	task.Tasks = cloneTasks(task.Tasks)
	return task
}

func cloneTasks(tasks Tasks) Tasks {
	if tasks == nil {
		return nil
	}
	cloned := make(Tasks, len(tasks))
	for i := range tasks {
		cloned[i] = tasks[i]
		cloned[i].Tasks = cloneTasks(tasks[i].Tasks)
	}
	return cloned
}

func pathWithIndex(path []int, index int) []int {
	next := make([]int, len(path)+1)
	copy(next, path)
	next[len(path)] = index
	return next
}

func IsCI() bool {
	return os.Getenv("CI") != "" ||
		os.Getenv("CONTINUOUS_INTEGRATION") != "" ||
		os.Getenv("BUILD_NUMBER") != "" ||
		os.Getenv("GITHUB_ACTIONS") != ""
}
