package parallel_task

import "context"

type TaskFunc func(ctx context.Context) error

type Task struct {
	Name string
	Func TaskFunc
}

type Result struct {
	Name       string
	Err        error
	PanicStack string
}
