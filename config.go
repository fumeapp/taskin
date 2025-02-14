package taskin

import (
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
)

type Config struct {
	Options         ConfigOptions
	Spinner         spinner.Spinner
	Colors          ConfigColors
	ProgressOptions []progress.Option
	Chars           ConfigChars
}

type ConfigOptions struct {
	ExitOnFailure bool
}

type ConfigColors struct {
	Spinner       lipgloss.Color
	Pending       lipgloss.Color
	Success       lipgloss.Color
	Failure       lipgloss.Color
	ParentStarted lipgloss.Color
}

type ConfigChars struct {
	ParentStarted string
	NotStarted    string
	Success       string
	Failure       string
}

var Defaults = Config{
	Options: ConfigOptions{
		ExitOnFailure: true,
	},
	Spinner: spinner.Dot,
	ProgressOptions: []progress.Option{
		progress.WithDefaultGradient(),
		progress.WithoutPercentage(),
		progress.WithWidth(10),
	},
	Colors: ConfigColors{
		Spinner:       lipgloss.Color("214"),
		Pending:       lipgloss.Color("lightgrey"),
		Success:       lipgloss.Color("46"),
		Failure:       lipgloss.Color("196"),
		ParentStarted: lipgloss.Color("214"),
	},
	Chars: ConfigChars{
		NotStarted:    "○", // Changed from "•" to "○"
		Success:       "✔",
		Failure:       "✗",
		ParentStarted: "»",
	},
}
