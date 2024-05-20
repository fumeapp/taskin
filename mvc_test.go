package taskin

import (
	"testing"
)

func TestModelInit(t *testing.T) {
	// Initialize a new Model
	m := &Model{}

	// Call the Init method
	cmd := m.Init()

	// Check if the returned command is not nil
	if cmd != nil {
		t.Errorf("Expected command to be not nil, got not nil")
	}
}

func TestModelUpdate(t *testing.T) {
	// Initialize a new Model
	m := &Model{}

	// Call the Update method with a dummy message
	newModel, cmd := m.Update("dummy message")

	// Check if the returned model is not nil
	if newModel == nil {
		t.Errorf("Expected model to be not nil, got nil")
	}

	// Check if the returned command is nil
	if cmd != nil {
		t.Errorf("Expected command to be nil, got non-nil")
	}
}

// Add more tests for other methods in the Model struct
