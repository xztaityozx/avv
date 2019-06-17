package pipeline

import (
	"context"
	"errors"
	"github.com/xztaityozx/avv/task"
)

type (
	PipeLine struct {
		Total  int
		Stages []Stage
	}

	Action func(ctx context.Context, t task.Task) (task.Task, error)

	Stage struct {
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
func (p *PipeLine) AddStage(w int, s chan task.Task, act Action) chan task.Task {
	st := Stage{
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
func (p PipeLine) Start(ctx context.Context, source []task.Task) error {

	return nil
}

// close close input, output, errorPipe
func (s Stage) close() {
	close(s.errorPipe)
	close(s.input)
	close(s.output)
}

// invoke start stage task
// returns:
//  - error:
func (s Stage) invoke(ctx context.Context) error {

	for i := 0; i < s.Worker; i++ {
		go func() {
			for v := range s.input {
				out, err := s.action(ctx, v)
				if err != nil {
					s.errorPipe <- err
					return
				}

				s.output <- out
			}
		}()
	}

	select {
	case <-ctx.Done():
		return errors.New("canceled")
	case err := <-s.errorPipe:
		return err
	}

}
