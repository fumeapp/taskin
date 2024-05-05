package taskin

import (
	"testing"
)

func TestTask(t *testing.T) {
	task := Task{
		Title: "Test Task",
		Task:  func(t *Task) error { return nil },
	}

	if task.Title != "Test Task" {
		t.Errorf("Expected task title to be 'Test Task', got '%s'", task.Title)
	}

	err := task.Task(&task)
	if err != nil {
		t.Errorf("Expected task function to return nil, got '%s'", err.Error())
	}
}

func TestRunner(t *testing.T) {
	runner := Runner{
		Task: Task{
			Title: "Test Task",
			Task:  func(t *Task) error { return nil },
		},
		State: NotStarted,
	}

	if runner.State != NotStarted {
		t.Errorf("Expected runner state to be 'NotStarted', got '%d'", runner.State)
	}

	if runner.Task.Title != "Test Task" {
		t.Errorf("Expected task title to be 'Test Task', got '%s'", runner.Task.Title)
	}

	err := runner.Task.Task(&runner.Task)
	if err != nil {
		t.Errorf("Expected task function to return nil, got '%s'", err.Error())
	}
}
