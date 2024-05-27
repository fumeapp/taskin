<h1 align="center">
    <img src="https://github.com/fumeapp/taskin/raw/main/taskin.png" width="300" />
 <br />
</h1>

<p align="center"><strong>📋 Add user-friendly tasks to your terminal </strong></p>
<br />

![Multi](/multi.gif)


[![Release](https://img.shields.io/github/v/release/fumeapp/taskin)](https://github.com/fumeapp/taskin/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/fumeapp/taskin)](https://goreportcard.com/report/github.com/fumeapp/taskin)
[![Go Reference](https://pkg.go.dev/badge/github.com/fumeapp/taskin.svg)](https://pkg.go.dev/github.com/fumeapp/taskin)
[![Lint](https://github.com/fumeapp/taskin/actions/workflows/lint.yml/badge.svg)](https://github.com/fumeapp/taskin/actions/workflows/lint.yml)
[![Tests](https://github.com/fumeapp/taskin/actions/workflows/test.yml/badge.svg)](https://github.com/fumeapp/taskin/actions/workflows/test.yml)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](https://github.com/fumeapp/taskin/pulls)
[![Examples](https://github.com/fumeapp/taskin/actions/workflows/examples.yml/badge.svg)](https://github.com/fumeapp/taskin/actions/workflows/examples.yml)
[![Phorm](https://img.shields.io/badge/Phorm-Ask-pink)](https://www.phorm.ai/query?projectId=4d6b35fb-2ee0-40a3-ad5f-952dc5f69365)


> [!TIP]
>
> All output is Github Action friendly! 
> You can view the output of each example [here](https://github.com/fumeapp/taskin/actions/workflows/examples.yml)



## Installation

```bash
go get github.com/fumeapp/taskin
```

## Examples

Simplest way to line up and fire off tasks

https://github.com/fumeapp/taskin/blob/3cd766c21e5eaba5edb33f38d3781d6cf814f9f9/examples/simple/main.go#L11-L33

![Simple](/simple.gif)


Using a progress bar for a task

https://github.com/fumeapp/taskin/blob/06b4d112f7d2dcf9fb4ee9b210f0be2d5cda03b5/examples/progress/main.go#L11-L24

![Progress](/progress.gif)


Customize colors, spinner, and progress bar

https://github.com/fumeapp/taskin/blob/3cd766c21e5eaba5edb33f38d3781d6cf814f9f9/examples/custom/main.go#L13-L54

![Custom](/custom.gif)



Nest tasks inside tasks

https://github.com/fumeapp/taskin/blob/3cd766c21e5eaba5edb33f38d3781d6cf814f9f9/examples/multi/main.go#L23-L34

![Multi](/multi.gif)



## Functionality in a task

### Change the title
Already demonstrated in most of the examples, you can change `t.Title` at any time

### Hide a view
Sometimes you might need to temporarily hide you task view in order to prompt a user for input.
You can do this by toggling the task.HideView boolean.

```go
Task: func(T *taskin.Task ) error {
	t.HideView = true
	if err := PromptForInput(); err != nil {
		t.HideView = false
        return err
    }
    t.HideView = false
	t.Title = "Input received"
	return nil
}

```
