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

func (d DummyTask) Run(ctx context.Context) Result {
	time.Sleep(time.Microsecond)
	return Result{
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

type SecondDummyTask struct {
	Task Task
}

func (s SecondDummyTask) Run(ctx context.Context) Result {
	time.Sleep(time.Microsecond)
	return Result{
		Status: true,
		Task:   s.Task,
	}
}

func (s SecondDummyTask) String() string {
	return ""
}

func (s SecondDummyTask) Self() Task {
	return s.Task
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
		as.Equal(Result{
			Task:   config.Default,
			Status: true,
		}, v)
	}

}

func TestPipeLine_Start(t *testing.T) {
	ctx := context.Background()
	p := PipeLine{}
	var in []ITask
	for i := 0; i < 10; i++ {
		in = append(in, DummyTask{Task: Task{}, Param: i})
	}

	ch := make(chan struct{})
	defer close(ch)

	as := assert.New(t)

	go func() {

		s, f, err := p.Start(ctx, in,
			Pipe{
				Name:       "first",
				Parallel:   5,
				RetryLimit: 2,
				Converter: func(task Task) ITask {
					return DummyTask{
						Task:  task,
						Param: 10,
					}
				},
			},
			Pipe{
				Name:       "second",
				Parallel:   8,
				RetryLimit: 2,
				Converter: func(task Task) ITask {
					return SecondDummyTask{
						Task: task,
					}
				},
			})
		as.Equal(10, len(s))
		as.Equal(0, len(f))
		as.Nil(err)
		ch <- struct{}{}
	}()

	select {
	case <-ctx.Done():
	case <-ch:
	}
}
