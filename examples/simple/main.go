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
			// sleep for 3 seconds then return nil
			Task: func(t *taskin.Task) error {
				for i := 0; i < 3; i++ {
					t.Title = fmt.Sprintf("Task 1: [%d/3] seconds have passed", i+1)
					time.Sleep(1 * time.Second)
				}
				return nil
			},
		},
		{
			Title: "Task 2",
			// sleep for 3 seconds then return nil
			Task: func(t *taskin.Task) error {
				for i := 0; i < 3; i++ {
					t.Title = fmt.Sprintf("Task 1: [%d/3] seconds have passed", i+1)
					time.Sleep(1 * time.Second)
				}
				return fmt.Errorf("task 2 failed")
			},
		},
	}, taskin.Defaults)
	err := tasks.Run()

	if err != nil {
		panic(err)
	}
}
