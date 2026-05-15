package taskin

import (
	"github.com/charmbracelet/bubbles/spinner"
	"testing"
)

func TestNewRunner(t *testing.T) {
	task := Task{
		Title: "Test Task",
		Task:  func(t *Task) error { return nil },
	}
	cfg := Config{
		Spinner: spinner.Dot,
	}

	runner := NewRunner(task, cfg)

	if runner.State != NotStarted {
		t.Errorf("Expected runner state to be 'NotStarted', got '%d'", runner.State)
	}

	if runner.Task.Title != "Test Task" {
		t.Errorf("Expected task title to be 'Test Task', got '%s'", runner.Task.Title)
	}
}

func TestRunnersRun(t *testing.T) {
	tasks := Tasks{
		Task{
			Title: "Test Task",
			Task:  func(t *Task) error { return nil },
		},
	}
	cfg := Config{
		Spinner: spinner.Dot,
	}

	runners := New(tasks, cfg)

	err := runners.Run()

	if err != nil {
		t.Errorf("Expected Run to return nil, got '%s'", err.Error())
	}
}

func TestNew(t *testing.T) {
	ran := false
	tasks := Tasks{
		Task{
			Title: "Test Task",
			Task: func(t *Task) error {
				ran = true
				return nil
			},
		},
	}
	cfg := Config{
		Spinner: spinner.Dot,
	}

	runners := New(tasks, cfg)

	if len(runners) != 1 {
		t.Errorf("Expected New to return 1 runner, got '%d'", len(runners))
	}

	if ran {
		t.Error("Expected New not to start tasks")
	}
}

func TestTaskProgress(t *testing.T) {
	task := Task{
		Title: "Test Task",
		Task:  func(t *Task) error { return nil },
	}
	cfg := Config{
		Spinner: spinner.Dot,
	}

	runner := NewRunner(task, cfg)
	runner.Task.Progress(1, 10)

	if runner.Task.ShowProgress.Current != 1 || runner.Task.ShowProgress.Total != 10 {
		t.Errorf("Expected TaskProgress to be 1/10, got %d/%d", runner.Task.ShowProgress.Current, runner.Task.ShowProgress.Total)
	}
}

func TestRunAppliesTaskUpdates(t *testing.T) {
	tasks := Tasks{
		{
			Title: "Test Task",
			Task: func(t *Task) error {
				t.SetTitle("Updated Task")
				t.Progress(2, 4)
				return nil
			},
		},
	}
	cfg := Defaults
	cfg.DisableUI = true

	runners := New(tasks, cfg)
	if err := runners.Run(); err != nil {
		t.Fatalf("Expected Run to return nil, got %s", err.Error())
	}

	if runners[0].State != Completed {
		t.Fatalf("Expected runner state to be completed, got %d", runners[0].State)
	}
	if runners[0].Task.Title != "Updated Task" {
		t.Fatalf("Expected task title update, got %q", runners[0].Task.Title)
	}
	if runners[0].Task.ShowProgress.Current != 2 || runners[0].Task.ShowProgress.Total != 4 {
		t.Fatalf("Expected progress 2/4, got %d/%d", runners[0].Task.ShowProgress.Current, runners[0].Task.ShowProgress.Total)
	}
}
