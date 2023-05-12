package parallel_task

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync"
)

func Execute(ctx context.Context, tasks []*Task) chan *Result {
	n := len(tasks)
	results := make(chan *Result, n)
	var wg sync.WaitGroup
	wg.Add(n)
	for _, task := range tasks {
		// in case unexpected cases occur
		tmpTask := task
		go func() {
			result := do(ctx, tmpTask)
			results <- result
			wg.Done()
		}()
	}

	wg.Wait()
	close(results)
	return results
}

func do(ctx context.Context, task *Task) *Result {
	result := &Result{
		Name: task.Name,
	}

	done := make(chan struct{}, 1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				result.PanicStack = string(debug.Stack())
				result.Err = fmt.Errorf("do task %v err: %v", task.Name, err)
				done <- struct{}{}
			}
		}()

		err := task.Func(ctx)
		if err != nil {
			result.Err = err
		}

		done <- struct{}{}
	}()

	select {
	case <-done:
	}

	return result
}
