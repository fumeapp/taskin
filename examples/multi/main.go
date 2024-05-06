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
					Title: "Pluck the polish and marans",
					Tasks: taskin.Tasks{
						{
							Title: "Pluck the polish",
							Task: func(t *taskin.Task) error {
								for i := 0; i < 3; i++ {
									t.Title = fmt.Sprintf(" [%d/3] polish plucked", i+1)
									time.Sleep(500 * time.Millisecond)
								}
								return nil
							},
						},
						{
							Title: "Pluck the marans",
							Task: func(t *taskin.Task) error {
								for i := 0; i < 3; i++ {
									t.Title = fmt.Sprintf(" [%d/3] marans plucked", i+1)
									time.Sleep(500 * time.Millisecond)
								}
								return nil
							},
						},
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
