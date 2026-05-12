package taskin

import (
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type TerminateWithError struct {
	Error error
}

type taskStartedMsg struct {
	Path []int
}

type taskUpdatedMsg struct {
	Path []int
	Task Task
}

type taskCompletedMsg struct {
	Path []int
	Task Task
}

type taskFailedMsg struct {
	Path  []int
	Task  Task
	Error error
}

type taskExecutionFinishedMsg struct{}

type TaskState int

const (
	NotStarted TaskState = iota
	Running
	Completed
	Failed
)

type Task struct {
	Title        string
	Task         func(*Task) error
	ShowProgress TaskProgress
	Bar          progress.Model
	Config       Config
	Tasks        Tasks
	HideView     bool
}

type TaskProgress struct {
	Current int
	Total   int
}

type Tasks []Task

type Runner struct {
	Task     Task
	State    TaskState
	Spinner  *spinner.Model
	Config   Config
	Children Runners
}

type Runners []Runner

type Model struct {
	Runners       Runners
	HideView      bool
	Shutdown      bool
	ShutdownError error
	taskMessages  chan tea.Msg
}
