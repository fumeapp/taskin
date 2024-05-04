package main

import (
    "github.com/charmbracelet/bubbles/spinner"
    "github.com/charmbracelet/lipgloss"
    tea "github.com/charmbracelet/bubbletea"
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
    Task  func(*Task) error
}

type Tasks []Task

type Runner struct {
    Task    Task
    State   TaskState
    Spinner spinner.Model
}

type Runners []Runner

func NewRunner(task Task) Runner {
    s := spinner.New(spinner.WithSpinner(spinner.Dot)) // Initialize with a spinner model
    s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205")) // Styling spinner
    return Runner{Task: task, State: NotStarted, Spinner: s}
}

func (r *Runners) Init() tea.Cmd {
    var cmds []tea.Cmd
    for i := range *r {
        // Here we directly use spinner.Tick, but since we're going to start tasks
        // and want them to run concurrently, we initialize those tasks elsewhere
        cmds = append(cmds, (*r)[i].Spinner.Tick)
    }
    // This function starts the tasks concurrently when the program initializes
    return tea.Batch(cmds...)
}

func (r *Runners) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmds []tea.Cmd

    switch msg := msg.(type) {

    case spinner.TickMsg:
        for i := range *r {
            if (*r)[i].State == Running || (*r)[i].State == NotStarted {
                // Update and capture new state of spinner and commands for the next tick
                newSpinner, cmd := (*r)[i].Spinner.Update(msg)
                (*r)[i].Spinner = newSpinner
                cmds = append(cmds, cmd)

                // Simulate task starting and state changing for demonstration
                if (*r)[i].State == NotStarted {
                    (*r)[i].State = Running
                }
            }
        }

        return r, tea.Batch(cmds...)
    }

    return r, nil
}

func (r *Runners) View() string {
    var view string
    for _, runner := range *r {
        status := ""
        switch runner.State {
        case NotStarted, Running:
            status = runner.Spinner.View() + " " + runner.Task.Title
        case Completed:
            status = "✔ " + runner.Task.Title
        case Failed:
            status = "✘ " + runner.Task.Title
        }
        view += lipgloss.NewStyle().Render(status) + "\n"
    }
    return view
}

func (r *Runners) Run() error {
    p := tea.NewProgram(r)
    _, err := p.Run()
    return err
}

func Listr(tasks Tasks) Runners {
    var runners Runners
    for _, task := range tasks {
        // Use NewRunner to ensure runners are initialized with spinners correctly
        runners = append(runners, NewRunner(task))
    }
    // simulate tasks running after the program starts
    go func() {
        time.Sleep(3 * time.Second) // Simulate delay before tasks change state
        for i := range runners {
            runners[i].State = Completed // Example of updating state, replace with real task logic
        }
    }()
    return runners
}

func main() {

    runners := Listr(Tasks{
        {
            Title: "Task 1",
            Task: func(t *Task) error {
                // The task logic should go here, replaced with sleep for demonstration
                time.Sleep(3 * time.Second)
                t.Title = "Task 1 Completed"
                return nil
            },
        },
    })
    err := runners.Run()

    if err != nil {
        panic(err)
    }
}
