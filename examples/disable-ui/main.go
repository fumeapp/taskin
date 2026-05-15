package main

import (
	"fmt"
	"time"

	"github.com/fumeapp/taskin"
)

func main() {
	tasks := taskin.Tasks{
		{
			Title: "Task with UI disabled",
			Task: func(t *taskin.Task) error {
				for i := 0; i < 3; i++ {
					t.SetTitle(fmt.Sprintf("Task with UI disabled: [%d/3] processing", i+1))
					time.Sleep(500 * time.Millisecond)
				}
				return nil
			},
		},
		{
			Title: "Parent task with children",
			Tasks: taskin.Tasks{
				{
					Title: "Child task 1",
					Task: func(t *taskin.Task) error {
						for i := 0; i < 2; i++ {
							t.SetTitle(fmt.Sprintf("Child task 1: [%d/2] working", i+1))
							time.Sleep(300 * time.Millisecond)
						}
						return nil
					},
				},
				{
					Title: "Child task 2",
					Task: func(t *taskin.Task) error {
						for i := 0; i < 2; i++ {
							t.SetTitle(fmt.Sprintf("Child task 2: [%d/2] working", i+1))
							time.Sleep(300 * time.Millisecond)
						}
						return nil
					},
				},
			},
		},
	}

	// Create configuration with UI disabled
	cfg := taskin.Defaults
	cfg.DisableUI = true

	runners := taskin.New(tasks, cfg)
	err := runners.Run()
	if err != nil {
		panic(err)
	}
}
