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
		Total      int
		Stages     []*Stage
		skip       bool
		Aggregator Aggregator
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

	Aggregator struct {
		name   string
		input  chan task.Task
		action func(ctx context.Context, box []task.Task) error
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
		skip:   true,
		Stages: []*Stage{},
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
	p.Stages = append(p.Stages, &st)

	return st.output
}

// Start start pipeline process
// params:
//  - ctx: context
//  - source: source tasks
// returns:
//  - error:
func (p *PipeLine) Start(ctx context.Context) error {

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
		go func(x *Stage) {
			defer wg.Done()
			err := x.invoke(ctx, bar)
			if err != nil {
				ch <- err
			}
		}(v)
	}

	wch := make(chan struct{})
	defer close(wch)

	go func() {
		pb.Wait()
		wg.Wait()
		wch <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return errors.New("canceled")
	case err := <-ch:
		return err
	case <-wch:
		return nil
	}
}

// invoke start stage task
// params:
//  - ctx: context
//  - bar: mpb.Bar
// returns:
//  - error:
func (s *Stage) invoke(ctx context.Context, bar *mpb.Bar) error {
	var wg sync.WaitGroup

	wch := make(chan struct{})
	defer close(wch)

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

	go func() {
		wg.Wait()
		close(s.output)
		wch <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return errors.New("canceled")
	case err := <-s.errorPipe:
		close(s.errorPipe)
		return err
	case <-wch:
		return nil
	}

}

// invoke start aggregate step
// params:
//  - ctx: context
//  - bar: mpb.Bar
// returns:
//  - error:
func (a *Aggregator) invoke(ctx context.Context, bar *mpb.Bar) error {

	// collect Task struct from input channel
	var box []task.Task
	for v := range a.input {
		box = append(box, v)
		bar.Increment()
	}

	// invoke aggregate func
	wch := make(chan error)
	defer close(wch)
	go func() {
		wch <- a.action(ctx, box)
	}()

	select {
	case <-ctx.Done():
		return errors.New("canceled")
	case err := <-wch:
		return err
	}
}

// AddAggregateStage add aggregator to pipeline
// params:
//  - in: source chan task.Task
//  - name: name of this aggregator
//  - act: something do in this stage
// returns:
//  - error:
func (p *PipeLine) AddAggregateStage(in chan task.Task, name string, act func(ctx context.Context, box []task.Task) error) error {

	if !p.skip {
		return errors.New("Already added aggregator\n")
	}

	p.skip = false
	p.Aggregator = Aggregator{
		input:  in,
		name:   name,
		action: act,
	}
	return nil
}
