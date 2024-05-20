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
			Title: "Task 2",
			Task: func(t *taskin.Task) error {
				for i := 0; i < 3; i++ {
					t.Title = fmt.Sprintf("Task 2: [%d/3] seconds have passed", i+1)
					time.Sleep(500 * time.Millisecond)
				}
				return nil
			},
		},
	}, taskin.Defaults)
	err := tasks.Run()

	if err != nil {
		panic(err)
	}
}
