package pipeline

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/xztaityozx/avv/parameters"
	"github.com/xztaityozx/avv/task"
	"math/rand"
	"testing"
	"time"
)

func TestNew(t *testing.T) {

	actual := New(10)
	expect := PipeLine{10, []*Stage{}, true, Aggregator{}}

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
	close(p.Stages[0].error)

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

func TestPipeLine_AddAggregateStage(t *testing.T) {
	in := make(chan task.Task)
	defer close(in)
	p := New(0)
	as := assert.New(t)

	name := "test-AddAggregatorStage"

	t.Run("NoError", func(t *testing.T) {
		err := p.AddAggregateStage(in, name, func(ctx context.Context, box []task.Task) error {
			if len(box)%2 == 0 {
				return nil
			}
			return errors.New("odd")
		})
		as.NoError(err)

	})

	t.Run("Error-Already-Added", func(t *testing.T) {
		err := p.AddAggregateStage(in, name, func(ctx context.Context, box []task.Task) error {
			return nil
		})
		as.Error(err)
	})

	as.Equal(in, p.Aggregator.input)
	as.Equal(name, p.Aggregator.name)

	t.Run("Action-returns", func(t *testing.T) {
		as.NoError(p.Aggregator.action(context.Background(), []task.Task{}))
		as.Error(p.Aggregator.action(context.Background(), []task.Task{{}}))
	})
}

func TestPipeLine_Start2(t *testing.T) {
	p := New(10)
	in := make(chan task.Task, 10)
	as := assert.New(t)

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

	err := p.AddAggregateStage(third, "agg", func(ctx context.Context, box []task.Task) error {
		var res []int

		for _, v := range box {
			res = append(res, v.Seed)
		}

		as.ElementsMatch([]int{400, 500, 600, 700, 800, 900, 1000, 1100, 1200, 1300}, res)
		return nil
	})

	as.NoError(err)

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

	select {
	case <-ctx.Done():
		as.Fail("canceled")
	case err := <-ech:
		as.NoError(err)
	}

}

func TestPipeLine_Start3(t *testing.T) {
	p := New(1)
	in := make(chan task.Task, 1)
	in <- task.Task{}
	close(in)

	_ = p.AddStage(1, in, "sleep", func(ctx context.Context, t task.Task) (i task.Task, e error) {

		ch := make(chan struct{})
		defer close(ch)

		go func() {
			time.Sleep(time.Second * 100000)
			ch<- struct{}{}
		}()

		select {
		case <-ctx.Done():
			return i, errors.New("canceled")
		case <-ch:
			return
		}
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ech := make(chan error)
	defer close(ech)
	go func() {
		ech <- p.Start(ctx)
	}()

	go func() {
		time.Sleep(time.Millisecond)
		cancel()
		logrus.Info("canceled")
	}()

	select {
	case err := <-ech:
		assert.Error(t, err)
	}
}
