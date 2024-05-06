package taskin

import (
	"os"
	"testing"
)

func TestRunners_Init(t *testing.T) {
	r := &Runners{
		// Initialize your Runners struct here
	}

	cmd := r.Init()

	// If Init is not implemented, it should return nil
	if cmd != nil {
		t.Errorf("Expected Init to return nil")
	}
}

func TestRunners_View(t *testing.T) {
	r := &Runners{
		// Initialize your Runners struct here
	}

	// Set the "CI" environment variable
	os.Setenv("CI", "true")

	view := r.View()

	// If "CI" is set and not all tasks are completed, View should return an empty string
	if view != "" {
		t.Errorf("Expected View to return an empty string")
	}
}
