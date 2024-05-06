<h1 align="center">
    <img src="https://github.com/fumeapp/taskin/raw/main/taskin.png" width="300" />
 <br />
</h1>

<p align="center"><strong>📋 Add user-friendly tasks to your CLI </strong></p>
<br />

![Multi](/multi.gif)


[![Release](https://img.shields.io/github/v/release/fumeapp/taskin)](https://github.com/fumeapp/taskin/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/fumeapp/taskin)](https://goreportcard.com/report/github.com/fumeapp/taskin)
[![Go Reference](https://pkg.go.dev/badge/github.com/fumeapp/taskin.svg)](https://pkg.go.dev/github.com/fumeapp/taskin)
[![Lint](https://github.com/fumeapp/taskin/actions/workflows/lint.yml/badge.svg)](https://github.com/fumeapp/taskin/actions/workflows/lint.yml)
[![Tests](https://github.com/fumeapp/taskin/actions/workflows/test.yml/badge.svg)](https://github.com/fumeapp/taskin/actions/workflows/test.yml)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](https://github.com/fumeapp/taskin/pulls)


## Installation

```bash
go get github.com/fumeapp/taskin
```

## Examples

Simplest way to line up and fire off tasks

![Simple](/simple.gif)

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
      Title: "Task 1",
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

Using a progress bar for a task

![Progress](/progress.gif)

https://github.com/fumeapp/taskin/blob/ad16df3b016a6fc2cb9c25102a7df97d079c9615/examples/progress/main.go

Customize colors, spinner, and progress bar

![Custom](/custom.gif)

```go
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
      Task: func(t *taskin.Task) error {
        for i := 0; i < 2; i++ {
          t.Title = fmt.Sprintf("Task 3 - [%d/%d]", i+1, 2)
          time.Sleep(1 * time.Second)
        }
        return nil
      },
    },
  }, taskin.Config{
    Spinner:        spinner.Moon,
    ProgressOption: progress.WithScaledGradient("#6667AB", "#34D399"),
  })
  err := tasks.Run()

  if err != nil {
    panic(err)
  }
}


```

Add multiple tasks to a single task

![Multi](/multi.gif)

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
      Title: "Mow the lawn",
      Task: func(t *taskin.Task) error {
        for i := 0; i < 3; i++ {
          t.Title = fmt.Sprintf("Mow the lawn: [%d/3] passes", i+1)
          time.Sleep(500 * time.Millisecond)
        }
        return nil
      },
    },
    {
      Title: "Pluck the Chickens",
      Tasks: taskin.Tasks{
        {
          Title: "Pluck the silkies",
          Task: func(t *taskin.Task) error {
            for i := 0; i < 3; i++ {
              t.Title = fmt.Sprintf(" [%d/3] silkies plucked", i+1)
              time.Sleep(500 * time.Millisecond)
            }
            return nil
          },
        },
        {
          Title: "Pluck the leghorns",
          Task: func(t *taskin.Task) error {
            for i := 0; i < 3; i++ {
              t.Title = fmt.Sprintf(" [%d/3] leghorns plucked", i+1)
              time.Sleep(500 * time.Millisecond)
            }
            return nil
          },
        },
      },
    },
    {
      Title: "Paint the house",
      Task: func(t *taskin.Task) error {
        for i := 0; i < 3; i++ {
          t.Progress(i+1, 5)
          t.Title = fmt.Sprintf("Paint the house: [%d/3] walls painted", i+1)
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

```
