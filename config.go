package taskin

import (
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
)

type Config struct {
	Spinner        spinner.Spinner
	Colors         ConfigColors
	ProgressOption progress.Option
}

type ConfigColors struct {
	Spinner lipgloss.Color
	Pending lipgloss.Color
	Success lipgloss.Color
	Failure lipgloss.Color
}

var Defaults = Config{
	Spinner:        spinner.Dot,
	ProgressOption: progress.WithDefaultGradient(),
	Colors: ConfigColors{
		Spinner: lipgloss.Color("214"),
		Pending: lipgloss.Color("21"),
		Success: lipgloss.Color("46"),
		Failure: lipgloss.Color("196"),
	},
}
