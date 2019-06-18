package pipeline

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/xztaityozx/avv/task"
	"testing"
)

func TestNew(t *testing.T) {
	actual := New(10)
	expect := PipeLine{10, []Stage{}}

	assert.Equal(t, expect, actual)
	assert.Equal(t, 10, actual.Total)
}

func TestPipeLine_AddStage(t *testing.T) {
	p := New(10)

	in := make(chan task.Task, 10)

	out := p.AddStage(10, in, "name", func(ctx context.Context, t task.Task) (i task.Task, e error) {
		return task.Task{}, nil
	})

	assert.Equal(t, 10, p.Stages[0].Worker)
	assert.Equal(t, in, p.Stages[0].input)
	assert.Equal(t, "name", p.Stages[0].name)

	close(out)
	close(p.Stages[0].input)
	close(p.Stages[0].errorPipe)

	_, err := p.Stages[0].action(context.Background(), task.Task{})
	assert.NoError(t, err)
}
