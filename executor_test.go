package parallel_task

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecute(t *testing.T) {
	tasks := make([]*Task, 0)
	tasks = append(tasks, &Task{
		Func: func(ctx context.Context) error {
			fmt.Println("task1")
			return nil
		},
		Name: "task1",
	})

	tasks = append(tasks, &Task{
		Func: func(ctx context.Context) error {
			fmt.Println("task2")
			return fmt.Errorf("err happened")
		},
		Name: "task2",
	})

	tasks = append(tasks, &Task{
		Func: func(ctx context.Context) error {
			fmt.Println("task3")
			panic("panic happened")
		},
		Name: "task3",
	})

	results := Execute(context.Background(), tasks)
	for result := range results {
		switch result.Name {
		case "task1":
			assert.Nil(t, result.Err)
		case "task2":
			assert.NotNil(t, result.Err)
		case "task3":
			assert.NotEqual(t, result.PanicStack, "")
		}
	}
}
