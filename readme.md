> we be taskin
```go
package main

import (
    "fmt"
    "github.com/fumeapp/taskin"
    "time"
)

func main() {

	tasks := taskin.Taskin(taskin.Tasks{
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
	})
	err := tasks.Run()

	if err != nil {
		panic(err)
	}
}
```

[Video](https://github.com/fumeapp/taskin/assets/967369/0305537b-f316-4ed3-80aa-0f654d2e818c)
