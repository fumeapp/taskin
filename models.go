package taskin

import (
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
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
}

type TaskProgress struct {
	Current int
	Total   int
}

type Tasks []Task

type Runner struct {
	Task    Task
	State   TaskState
	Spinner spinner.Model
	Config  Config
}

type Runners []Runner

type Config struct {
	Spinner spinner.Spinner
	Colors  ConfigColors
}

type ConfigColors struct {
	Spinner lipgloss.Color
	Pending lipgloss.Color
	Success lipgloss.Color
	Failure lipgloss.Color
}

var Defaults = Config{
	Spinner: spinner.Dot,
	Colors: ConfigColors{
		Spinner: lipgloss.Color("214"),
		Pending: lipgloss.Color("21"),
		Success: lipgloss.Color("46"),
		Failure: lipgloss.Color("196"),
	},
}
