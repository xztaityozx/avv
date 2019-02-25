package cmd

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewDispatcher(t *testing.T) {
	as := assert.New(t)

	expect := Dispatcher{}
	as.Equal(expect, NewDispatcher(""))
}

type DummyTask struct {
	Task  Task
	Param int
}

func (d DummyTask) Run(ctx context.Context) TaskResult {
	time.Sleep(time.Microsecond)
	return TaskResult{
		Status: true,
		Task:   d.Task,
	}
}

func (d DummyTask) String() string {
	return fmt.Sprint(d.Param)
}

func (d DummyTask) Self() Task {
	return d.Task
}


func TestDispatcher_Dispatch(t *testing.T) {
	d := NewDispatcher("")
	ctx, can := context.WithCancel(context.Background())
	defer can()

	var tasks []ITask
	for i := 0; i < 20; i++ {
		tasks = append(tasks, DummyTask{Task: config.Default, Param: i})
	}

	res := d.Dispatch(ctx, 4, tasks)

	as := assert.New(t)
	as.Equal(20, d.Size)
	as.Equal(20, len(res))

	for _, v := range res {
		as.Equal(TaskResult{
			Task:   config.Default,
			Status: true,
		}, v)
	}

}

