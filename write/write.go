package write

import (
	"context"
	"github.com/xztaityozx/avv/parameters"
	"github.com/xztaityozx/avv/task"
)

type Write struct {
	Tmp parameters.Templates
}

// Invoke start write files for simulation(IStage)
// params:
//  - ctx: context
//  - t: task
// returns:
//  - Task:
//  - error:
func (w Write) Invoke(ctx context.Context, t task.Task) (task.Task, error) {
	err := t.MakeFiles(w.Tmp)
	return t, err
}
