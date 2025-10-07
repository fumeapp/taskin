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
		Options: ConfigOptions{
			ExitOnFailure: true,
		},
		Spinner: spinner.Dot,
		Colors: ConfigColors{
			Spinner:       lipgloss.Color("214"),
			Pending:       lipgloss.Color("21"),
			Success:       lipgloss.Color("46"),
			Failure:       lipgloss.Color("196"),
			ParentStarted: lipgloss.Color("214"),
		},
		Chars: ConfigChars{
			ParentStarted: "»",
			NotStarted:    "•",
			Success:       "✔",
			Failure:       "✘",
		},
		ProgressOptions: []progress.Option{
			progress.WithDefaultGradient(),
			progress.WithoutPercentage(),
			progress.WithWidth(10),
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

	if config.Options.ExitOnFailure != true {
		t.Errorf("Expected ExitOnFailure to be 'true', got '%v'", config.Options.ExitOnFailure)
	}
}

func TestDisableUI(t *testing.T) {
	tasks := Tasks{
		{
			Title: "Test task",
			Task: func(task *Task) error {
				task.Title = "Test task completed"
				return nil
			},
		},
	}

	// Test with UI disabled
	cfg := Defaults
	cfg.DisableUI = true

	runners := New(tasks, cfg)
	
	// Verify that spinners are not initialized when UI is disabled
	for _, runner := range runners {
		if runner.Spinner != nil {
			t.Error("Expected spinner to be nil when DisableUI is true")
		}
	}

	err := runners.Run()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}
