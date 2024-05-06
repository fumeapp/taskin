package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/fumeapp/taskin"
	"time"
)

func main() {

	tasks := taskin.New(taskin.Tasks{
		{
			Title: "Task 1",
			// sleep for 3 seconds then return nil
			Task: func(t *taskin.Task) error {
				for i := 0; i < 2; i++ {
					t.Title = fmt.Sprintf("Task 1 - [%d/%d]", i+1, 2)
					time.Sleep(1 * time.Second)
				}
				return nil
			},
		},
		{
			Title: "Task 2 Progress",
			// sleep for 3 seconds then return nil
			Task: func(t *taskin.Task) error {
				for i := 0; i < 5; i++ {
					t.Progress(i+1, 5)
					t.Title = fmt.Sprintf("Task 2 - [%d/%d]", i+1, 5)
					time.Sleep(1 * time.Second)
				}
				return nil
			},
		},
		{
			Title: "Task 3",
			// sleep for 3 seconds then return nil
			Task: func(t *taskin.Task) error {
				for i := 0; i < 2; i++ {
					t.Title = fmt.Sprintf("Task 3 - [%d/%d]", i+1, 2)
					time.Sleep(1 * time.Second)
				}
				return nil
			},
		},
	}, taskin.Config{
		Spinner: spinner.Moon,
		ProgressOptions: []progress.Option{
			progress.WithScaledGradient("#6667AB", "#34D399"),
		},
	})
	err := tasks.Run()

	if err != nil {
		panic(err)
	}
}
