package pipeline

import (
	"context"
	"errors"
	"github.com/fatih/color"
	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
	"github.com/xztaityozx/avv/task"
	"sync"
)

type (
	PipeLine struct {
		Total  int
		Stages []Stage
	}

	Action func(ctx context.Context, t task.Task) (task.Task, error)

	Stage struct {
		name      string
		Worker    int
		input     chan task.Task
		output    chan task.Task
		errorPipe chan error
		action    Action
	}
)

// New make struct PipeLine struct
// params:
//  - t: total of data
// returns:
//  - PipeLine:
func New(t int) PipeLine {
	return PipeLine{
		Total:  t,
		Stages: []Stage{},
	}
}

// AddStage add stage to pipeline
// params:
//  - w: workers
//  - s: source chan
//  - act: function for this stage
// returns:
//  - chan task.Task: output chan
func (p *PipeLine) AddStage(w int, s chan task.Task, name string, act Action) chan task.Task {
	st := Stage{
		name:      name,
		action:    act,
		input:     s,
		output:    make(chan task.Task, p.Total),
		errorPipe: make(chan error),
		Worker:    w,
	}
	p.Stages = append(p.Stages, st)

	return st.output
}

// Start start pipeline process
// params:
//  - ctx: context
//  - source: source tasks
// returns:
//  - error:
func (p PipeLine) Start(ctx context.Context) error {
	var wg sync.WaitGroup
	wg.Add(len(p.Stages))

	pb := mpb.NewWithContext(ctx)

	workingMSG := color.New(color.FgHiYellow).Sprint("processing...")
	finishMSG := color.New(color.FgHiGreen).Sprint("done!")

	ch := make(chan error)
	defer close(ch)

	for _, v := range p.Stages {
		barName := color.New(color.FgHiCyan).Sprint(v.name)
		bar := pb.AddBar(int64(p.Total),
			mpb.BarStyle("┃██▒┃"),
			mpb.BarWidth(50),
			mpb.PrependDecorators(
				decor.Name(barName, decor.WC{W: len(barName) + 1, C: decor.DidentRight}),
			),
			mpb.AppendDecorators(
				decor.Name("   "),
				decor.Percentage(decor.WC{W: 5}),
				decor.Name(" | "),
				decor.CountersNoUnit("%d / %d", decor.WCSyncWidth),
				decor.Name(" | "),
				decor.OnComplete(decor.Name(workingMSG), finishMSG)))
		go func() {
			defer wg.Done()
			err := v.invoke(ctx, bar)
			if err != nil {
				ch <- err
			}
		}()
	}

	pb.Wait()
	wg.Done()

	c := func() {
		for _, v := range p.Stages {
			v.close()
		}
	}

	select {
	case <-ctx.Done():
		c()
		return errors.New("canceled")
	case err := <-ch:
		c()
		return err
	}
}

// close close input, output, errorPipe
func (s Stage) close() {
	close(s.errorPipe)
	close(s.input)
	close(s.output)
}

// invoke start stage task
// params:
//  - ctx: context
//  - bar: mpb.Bar
// returns:
//  - error:
func (s Stage) invoke(ctx context.Context, bar *mpb.Bar) error {
	var wg sync.WaitGroup

	for i := 0; i < s.Worker; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for v := range s.input {
				out, err := s.action(ctx, v)
				if err != nil {
					s.errorPipe <- err
					return
				}

				s.output <- out
				bar.Increment()
			}
		}()
	}

	wg.Wait()

	select {
	case <-ctx.Done():
		return errors.New("canceled")
	case err := <-s.errorPipe:
		return err
	}

}
