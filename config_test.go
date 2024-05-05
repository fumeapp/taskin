package taskin

import (
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
	"reflect"
	"testing"
)

func TestConfig(t *testing.T) {
	config := Config{
		Spinner:        spinner.Dot,
		ProgressOption: progress.WithDefaultGradient(),
		Colors: ConfigColors{
			Spinner: lipgloss.Color("214"),
			Pending: lipgloss.Color("21"),
			Success: lipgloss.Color("46"),
			Failure: lipgloss.Color("196"),
		},
	}

	if !reflect.DeepEqual(config.Spinner.Frames, spinner.Dot.Frames) {
		t.Errorf("Expected spinner frames to be equal to 'Dot' frames")
	}

	if config.Colors.Spinner != lipgloss.Color("214") {
		t.Errorf("Expected spinner color to be '214', got '%s'", config.Colors.Spinner)
	}

	if config.Colors.Pending != lipgloss.Color("21") {
		t.Errorf("Expected pending color to be '21', got '%s'", config.Colors.Pending)
	}

	if config.Colors.Success != lipgloss.Color("46") {
		t.Errorf("Expected success color to be '46', got '%s'", config.Colors.Success)
	}

	if config.Colors.Failure != lipgloss.Color("196") {
		t.Errorf("Expected failure color to be '196', got '%s'", config.Colors.Failure)
	}
}
