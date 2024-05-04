package main

import (
	"time"
)

type TaskState int

const (
	NotStarted TaskState = iota
	Running
	Completed
	Failed
)

type Task struct {
	Title string
	Task  func() error
}
type Tasks []Task

type Runner struct {
	Task  Task
	State TaskState
}
type Runners []Runner

func (r *Runners) Run() error {
	for _, runner := range *r {
		runner.State = Running
		err := runner.Task.Task()
		if err != nil {
			runner.State = Failed
			return err
		}
		runner.State = Completed
	}
	return nil

}

func Listr(tasks Tasks) Runners {

	var runners Runners

	for _, task := range tasks {
		runners = append(runners, Runner{Task: task, State: NotStarted})
	}

	return runners
}

func main() {

	runners := Listr(Tasks{
		{
			Title: "Task 1",
			// sleep for 3 seconds then return nil
			Task: func() error {
				time.Sleep(2 * time.Second)
				return nil
			},
		},
	})
	err := runners.Run()

	if err != nil {
		panic(err)
	}
}
