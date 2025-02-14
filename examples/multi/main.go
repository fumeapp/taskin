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
					Title: "Pluck the silkies [0/3]",
					Task: func(t *taskin.Task) error {
						for i := 0; i < 3; i++ {
							t.Title = fmt.Sprintf("Pluck the silkies [%d/3]", i+1)
							time.Sleep(500 * time.Millisecond)
						}
						return nil
					},
				},

				{
					Title: "Pluck the Polish and Marans",
					Tasks: taskin.Tasks{
						{
							Title: "[0/3] Pluck the Polish",
							Task: func(t *taskin.Task) error {
								for i := 0; i < 3; i++ {
									t.Title = fmt.Sprintf("[%d/3] Pluck the Polish", i+1)
									time.Sleep(500 * time.Millisecond)
								}
								return nil
							},
						},
						{
							Title: "[0/3] Pluck the Marans",
							Task: func(t *taskin.Task) error {
								for i := 0; i < 3; i++ {
									t.Title = fmt.Sprintf("[%d/3] Pluck the Marans", i+1)
									time.Sleep(500 * time.Millisecond)
								}
								return nil
							},
						},
					},
				},

				{
					Title: "[0/3] Pluck the leghorns",
					Task: func(t *taskin.Task) error {
						for i := 0; i < 3; i++ {
							t.Title = fmt.Sprintf("[%d/3] Pluck the Leghorns", i+1)
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
					t.Progress(i+1, 3)
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
