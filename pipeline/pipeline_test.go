package pipeline

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/xztaityozx/avv/parameters"
	"github.com/xztaityozx/avv/task"
	"math/rand"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	actual := New(10)
	expect := PipeLine{10, []*Stage{}, false, Aggregator{}}

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

func TestPipeLine_Start(t *testing.T) {
	p := New(10)
	in := make(chan task.Task, 10)

	first := p.AddStage(2, in, "+4 and sleep random ms", func(ctx context.Context, t task.Task) (i task.Task, e error) {
		i = t
		i.Seed += 4

		x := rand.New(rand.NewSource(time.Now().Unix()))
		time.Sleep(time.Millisecond*time.Duration(x.Int()%1000) + time.Millisecond*100)

		return
	})

	second := p.AddStage(3, first, "*10 and sleep random ms", func(ctx context.Context, t task.Task) (i task.Task, e error) {
		i = t
		i.Parameters.Seed *= 10
		x := rand.New(rand.NewSource(time.Now().Unix()))
		time.Sleep(time.Millisecond*time.Duration(x.Int()%1000) + time.Millisecond*100)

		return
	})

	third := p.AddStage(5, second, "*10", func(ctx context.Context, t task.Task) (i task.Task, e error) {
		i = t
		i.Seed *= 10
		return
	})

	ech := make(chan error)
	defer close(ech)

	ctx := context.Background()

	go func() {
		for i := 0; i < 10; i++ {
			in <- task.Task{Parameters: parameters.Parameters{
				Seed: i,
			}}
		}
		close(in)
	}()

	go func() {
		ech <- p.Start(ctx)
	}()

	as := assert.New(t)

	select {
	case <-ctx.Done():
		as.Fail("canceled")
	case err := <-ech:
		as.NoError(err)
	}

	var res []int
	for v := range third {
		res = append(res, v.Seed)
	}

	as.ElementsMatch([]int{400, 500, 600, 700, 800, 900, 1000, 1100, 1200, 1300}, res)
}
