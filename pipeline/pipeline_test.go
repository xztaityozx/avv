package pipeline

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/xztaityozx/avv/parameters"
	"github.com/xztaityozx/avv/task"
	"testing"
)

type d struct {
	act func(task task.Task) task.Task
}

func (d d) Invoke(ctx context.Context, t task.Task) (task.Task, error) {
	return d.act(t), nil
}

func TestNew(t *testing.T) {
	actual := New(1, 2)
	assert.Equal(t, 1, actual.Total)
	assert.Equal(t, 2, actual.maxRetry)
	assert.Equal(t, 0, len(actual.stages))
}

func TestPipeLine_AddStage(t *testing.T) {
	p := New(1, 2)

	s := make(chan task.Task, 1)
	s <- task.Task{}
	close(s)

	_ = p.AddStage(1, s, "test", d{})
	actual := p.stages[0]
	as := assert.New(t)
	as.Equal(1, actual.Worker)
	as.Equal("test", actual.name)
	as.Equal(d{}, actual.iStage)
}

func TestPipeLine_Start(t *testing.T) {
	size := 10
	in := make(chan task.Task, size)
	for i := 0; i < size; i++ {
		in <- task.Task{Parameters: parameters.Parameters{Seed: i}}
	}
	close(in)

	p := New(size, 4)
	first := p.AddStage(2, in, "first", d{act: func(task task.Task) task.Task {
		return task
	}})
	second := p.AddStage(2, first, "second", d{act: func(task task.Task) task.Task {
		task.Seed += 10
		return task
	}})
	third := p.AddStage(2, second, "third", d{act: func(task task.Task) task.Task {
		task.Seed *= 10
		return task
	}})

	as := assert.New(t)
	ctx := context.Background()
	ch := make(chan error, 1)
	defer close(ch)
	go func() {
		ch <- p.Start(ctx)
	}()

	select {
	case <-ctx.Done():
		as.Fail("canceled")
	case err := <-ch:
		as.NoError(err)
	}

	var box []int
	for v := range third {
		box = append(box, v.Seed)
	}

	as.ElementsMatch([]int{100, 110, 120, 130, 140, 150, 160, 170, 180, 190}, box)
}
