package main

import (
	"fmt"
	"context"
)

var TaskHandlers = map[TaskType]TaskHandler {
	"default": TaskUnrecognizedHandler,
}

type TaskHandler func(ctx context.Context, task *Task)

func TaskUnrecognizedHandler(ctx context.Context, task *Task) {
	fmt.Println("Task unrecognized:", task)
}

type TaskType string

type Task struct {
	typ TaskType
	master string
	run string
}

func (t Task) String() string {
	return fmt.Sprintf("Task %s[M:%s, R:%s]", t.typ, t.master, t.run)
}

