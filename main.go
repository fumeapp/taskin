package main

import (
 "github.com/charmbracelet/bubbles/spinner"
 "github.com/charmbracelet/lipgloss"
 "time"

 tea "github.com/charmbracelet/bubbletea"
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
 Task  Task
 State TaskState
 Msg   tea.Msg
 View  string
}

type Runners []Runner

func (r *Runners) Init() tea.Cmd {
 return tea.Batch(func() tea.Cmd {
  for i := range *r {
   (*r)[i].State = Running
   return tea.Tick(time.Second, func(t time.Time) tea.Msg {
    return &(*r)[i]
   })
  }
  return nil
 }())
}

func (r *Runners) Run() error {
 p := tea.NewProgram(r)
 _, err :=  p.Run()
 return err
}

func (r *Runners) View() string {
 var view string
 s := spinner.New()
 s.Spinner = spinner.Dot

 for _, runner := range *r {
  switch runner.State {
  case Running:
   view += lipgloss.NewStyle().Render(" " + s.View() + " " + runner.Task.Title)
  case Completed:
   view += lipgloss.NewStyle().Render(" ✔ " + runner.Task.Title)
  case Failed:
   view += lipgloss.NewStyle().Render(" ✘ " + runner.Task.Title)
  }
  view += "\n"
 }
 return view
}

func (r *Runners) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
 allDone := true
 s := spinner.NewModel()
 s.Tick() // Update the spinner's frame

 for i := range *r {
  if msg, ok := msg.(*Runner); ok && msg == &(*r)[i] {
   err := (*r)[i].Task.Task(&(*r)[i].Task)
   if err != nil {
    (*r)[i].State = Failed
   } else {
    (*r)[i].State = Completed
   }
  }

  // If any task is not completed or failed, set allDone to false
  if (*r)[i].State != Completed && (*r)[i].State != Failed {
   allDone = false
  }
 }

 // If all tasks are done, return a Quit command
 if allDone {
  return r, tea.Quit
 }

 return r, nil
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
   Task: func(t *Task) error {
    time.Sleep(1 * time.Second)
    t.Title = "1 second has passed"
    time.Sleep(1 * time.Second)
    t.Title = "2 seconds have passed"
    time.Sleep(1 * time.Second)
    return nil
   },
  },
 })
 err := runners.Run()

 if err != nil {
  panic(err)
 }
}
