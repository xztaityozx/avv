package pipeline

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
	"github.com/xztaityozx/avv/task"
	"golang.org/x/xerrors"
	"strings"
	"sync"
)

type (
	PipeLine struct {
		Total    int
		stages   []*Stage
		maxRetry int
	}

	Stage struct {
		name   string
		Worker int
		input  chan task.Task
		output chan task.Task
		error  chan error
		iStage IStage
	}

	IStage interface {
		Invoke(ctx context.Context, t task.Task) (task.Task, error)
	}
)

// New make struct PipeLine struct
// params:
//  - t: Total of data
// returns:
//  - PipeLine:
func New(t, m int) PipeLine {
	return PipeLine{
		Total:    t,
		stages:   []*Stage{},
		maxRetry: m,
	}
}

// AddStage add stage to pipeline
// params:
//  - w: number of workers
//  - s: source chan
//  - name: name of this stage
//  - is: struct that implemented IStage interface
func (p *PipeLine) AddStage(w int, s chan task.Task, name string, is IStage) chan task.Task {
	st := Stage{
		name:   name,
		input:  s,
		output: make(chan task.Task, p.Total),
		error:  make(chan error, p.Total),
		Worker: w,
		iStage: is,
	}

	p.stages = append(p.stages, &st)
	return st.output
}

// Start start pipeline process
// params:
//  - ctx: context
//  - source: source tasks
// returns:
//  - error:
func (p *PipeLine) Start(ctx context.Context) error {

	// padding name
	max := 0
	for _, v := range p.stages {
		if max < len(v.name) {
			max = len(v.name)
		}
	}

	for i := range p.stages {
		l := max - len(p.stages[i].name)
		if l == 0 {
			continue
		}
		p.stages[i].name += strings.Repeat(" ", l)
	}

	// make WaitGroup
	var wg sync.WaitGroup
	// add size of stages to wg
	wg.Add(len(p.stages))

	// make Progressbar
	pb := mpb.NewWithContext(ctx)

	workingMSG := color.New(color.FgHiYellow).Sprint("processing...")
	finishMSG := color.New(color.FgHiGreen).Sprint("done!")

	ch := make(chan error)
	defer close(ch)

	// Start stages
	for _, v := range p.stages {
		// make progressbar
		barName := color.New(color.FgHiCyan).Sprint(v.name)
		bar := makeBar(p.Total, barName, workingMSG, finishMSG, pb)

		// start stage
		go func(x *Stage) {
			defer wg.Done()
			err := x.invoke(ctx, bar, p.maxRetry)
			if err != nil {
				ch <- xerrors.Errorf("Failed Stage %s : %s", x.name, err)
			}
		}(v)
	}

	wch := make(chan struct{})

	go func() {
		pb.Wait()
		wg.Wait()
		wch <- struct{}{}
	}()

	select {
	case err := <-ch:
		return err
	case <-wch:
		close(wch)
		return nil
	}
}

// invoke start stage task
// params:
//  - ctx: context
//  - bar: mpb.Bar
//  - max: limit of retry
// returns:
//  - error:
func (s *Stage) invoke(ctx context.Context, bar *mpb.Bar, max int) error {
	var wg sync.WaitGroup

	for i := 0; i < s.Worker; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for v := range s.input {

				var err error
				var out task.Task
				for i := 0; i < max; i++ {
					out, err = s.iStage.Invoke(ctx, v)
					// continue retry
					if err == nil {
						break
					}
				}

				if err != nil {
					s.error <- err
				} else {
					s.output <- out
				}
				bar.Increment()
			}
		}()
	}

	wch := make(chan struct{})
	defer close(wch)

	go func() {
		wg.Wait()
		close(s.output)
		wch <- struct{}{}
	}()

	select {
	case <-wch:
		close(s.error)
		for err := range s.error {
			if err != nil {
				return err
			}
		}
		return nil
	}

}

func makeBar(total int, barName, workingMSG, finishMSG string, pb *mpb.Progress) *mpb.Bar {

	side := color.New(color.FgHiGreen).Sprint(string('\u258D'))
	done := color.New(color.FgCyan).Sprint(string('\u2588'))
	now := color.New(color.FgHiBlue).Sprint(string('\u2588'))
	wait := fmt.Sprint(string('\u2591'))

	style := side + done + now + wait + side

	return pb.AddBar(int64(total),
		mpb.BarStyle(style),
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
}
