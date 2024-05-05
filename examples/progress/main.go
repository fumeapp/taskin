package main

import (
	"fmt"
	"github.com/fumeapp/taskin"
	"time"
)

func main() {

	tasks := taskin.New(taskin.Tasks{
		{
			Title: "Progress Task",
			Task: func(t *taskin.Task) error {
				for i := 0; i < 5; i++ {
					t.Progress(i+1, 5)
					t.Title = fmt.Sprintf("Task 2 - [%d/%d]", i+1, 5)
					time.Sleep(1 * time.Second)
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
