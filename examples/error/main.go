package main

import (
	"fmt"
	"github.com/fumeapp/taskin"
	"time"
)

func main() {

	tasks := taskin.New(taskin.Tasks{
		{
			Title: "Task 1",
			Task: func(t *taskin.Task) error {
				for i := 0; i < 3; i++ {
					t.Title = fmt.Sprintf("Task 1: [%d/3] seconds have passed", i+1)
					time.Sleep(500 * time.Millisecond)
				}
				return nil
			},
		},
		{
			Title: "Task 2: Error",
			Task: func(t *taskin.Task) error {
				return fmt.Errorf("task 2 failed")
			},
		},
		{
			Title: "Task 3",
			Task: func(t *taskin.Task) error {
				for i := 0; i < 3; i++ {
					t.Title = fmt.Sprintf("Task 3: [%d/3] seconds have passed", i+1)
					time.Sleep(500 * time.Millisecond)
				}
				return nil
			},
		},
	}, taskin.Defaults)
	// }, taskin.Config{Options: taskin.ConfigOptions{ExitOnFailure: false}})
	err := tasks.Run()

	if err != nil {
		panic(err)
	}
}
