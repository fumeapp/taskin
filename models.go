package taskin

import (
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
)

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
