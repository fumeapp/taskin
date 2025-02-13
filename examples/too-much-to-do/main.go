package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fumeapp/taskin"
)

type Library struct {
	Name string
}

func main() {
	libraries := generateLibraries(20)
	runners := NewLibraryDownloader(libraries)

	runners.Run()
}

func generateLibraries(count int) []Library {
	libraries := make([]Library, count)
	for i := 0; i < count; i++ {
		libraries[i] = Library{
			Name: fmt.Sprintf("Library%d", i+1),
		}
	}
	return libraries
}

func NewLibraryDownloader(libraries []Library) taskin.Runners {
	var tasks []taskin.Task

	for _, library := range libraries {
		tasks = append(tasks, createTasksForLibrary(library))
	}

	if len(tasks) == 0 {
		fmt.Println("No tasks to run")
		return nil
	}

	return taskin.New(tasks, taskin.Defaults)
}

func createTasksForLibrary(library Library) taskin.Task {
	return taskin.Task{
		Title: fmt.Sprintf("Process %s", library.Name),
		Tasks: taskin.Tasks{
			{
				Title: "Download",
				Task:  createSimulatedTask(500, 2000),
			},
			{
				Title: "Unarchive",
				Task:  createSimulatedTask(200, 1000),
			},
			{
				Title: "Process",
				Task:  createSimulatedTask(300, 1500),
			},
		},
	}
}

func createSimulatedTask(minDuration, maxDuration int) func(*taskin.Task) error {
	return func(_ *taskin.Task) error {
		duration := time.Duration(rand.Intn(maxDuration-minDuration)+minDuration) * time.Millisecond
		time.Sleep(duration)

		return nil
	}
}
