package taskin

import (
	"github.com/charmbracelet/bubbles/spinner"
	"testing"
)

func TestRunners_Init(t *testing.T) {
	runners := &Runners{
		NewRunner(Task{Title: "Test Task", Task: func(t *Task) error { return nil }}, Config{}),
	}

	cmd := runners.Init()

	if cmd == nil {
		t.Errorf("Expected Init to return a non-nil Cmd")
	}
}

func TestRunners_Update(t *testing.T) {
	runners := &Runners{
		// Initialize with some test data
	}

	model, cmd := runners.Update(spinner.TickMsg{})

	if model == nil {
		t.Errorf("Expected Update to return a non-nil Model")
	}

	if cmd == nil {
		t.Errorf("Expected Update to return a non-nil Cmd")
	}
}

func TestRunners_View(t *testing.T) {
	runners := &Runners{
		NewRunner(Task{Title: "Test Task", Task: func(t *Task) error { return nil }}, Config{}),
	}

	view := runners.View()

	if view == "" {
		t.Errorf("Expected View to return a non-empty string")
	}
}
