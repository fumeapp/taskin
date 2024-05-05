
# Taskin

<p align="center">
  <img src="https://raw.githubusercontent.com/fumeapp/taskin/857a1b8cbeda577a751b5c7f38885995a894169f/taskin.png" width="300" />
</p>

![Simple](/simple.gif)

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
          t.Title = fmt.Sprintf("Task 2: [%d/3] seconds have passed", i+1)
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

![Progress](/progress.gif)


```go
package main

import (
	"fmt"
	"github.com/fumeapp/taskin"
	"time"
)

func main() {

	tasks := taskin.New(taskin.Tasks{
		{
			Title: "Progress",
			Task: func(t *taskin.Task) error {
				for i := 0; i < 5; i++ {
					t.Progress(i+1, 5)
					t.Title = fmt.Sprintf("Progress [%d/%d]", i+1, 5)
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
```
