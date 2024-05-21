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

var DefaultConfigColors = ConfigColors{
	Spinner:       lipgloss.Color("214"),
	Pending:       lipgloss.Color("lightgrey"),
	Success:       lipgloss.Color("46"),
	Failure:       lipgloss.Color("196"),
	ParentStarted: lipgloss.Color("214"),
}

var AnsiColors = ConfigColors{
	Spinner:       lipgloss.Color("3"), // Yellow, similar to "214"
	Pending:       lipgloss.Color("7"), // White, similar to "lightgrey"
	Success:       lipgloss.Color("2"), // Green, similar to "46"
	Failure:       lipgloss.Color("1"), // Red, similar to "196"
	ParentStarted: lipgloss.Color("3"), // Yellow, similar to "214"
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
	Colors: DefaultConfigColors,
	Chars: ConfigChars{
		ParentStarted: "»",
		NotStarted:    "•",
		Success:       "✔",
		Failure:       "✘",
	},
}

var AnsiDefaults = Config{
	Options: ConfigOptions{
		ExitOnFailure: true,
	},
	Spinner: spinner.Dot,
	ProgressOptions: []progress.Option{
		progress.WithDefaultGradient(),
		progress.WithoutPercentage(),
		progress.WithWidth(10),
	},
	Colors: AnsiColors,
	Chars: ConfigChars{
		ParentStarted: "»",
		NotStarted:    "•",
		Success:       "✔",
		Failure:       "✘",
	},
}
