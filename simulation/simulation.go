package simulation

import (
	"context"
	"github.com/xztaityozx/avv/extract"
	"github.com/xztaityozx/avv/task"
)

type Simulation struct {
	HSPICE HSPICE
	WaveView extract.WaveView
}

func (s Simulation) Invoke(ctx context.Context, t task.Task) (task.Task, error) {
	simResult, err := s.HSPICE.Invoke(ctx, t)
	if err != nil {
		return task.Task{}, err
	}
	return s.WaveView.Invoke(ctx, simResult)
}

