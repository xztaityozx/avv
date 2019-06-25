package pipeline

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/xztaityozx/avv/task"
	"testing"
)

type d struct {
}

func (d) Invoke(ctx context.Context, t task.Task) (task.Task, error) {
	return t, nil
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
