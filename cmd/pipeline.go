package cmd

import (
	"context"
	"github.com/sirupsen/logrus"
	"time"
)

type (
	PipeLine struct {
		Tasks []ITask
	}
)

func NewPipeLine(tasks []ITask) PipeLine {
	var x []ITask

	for _, r := range tasks {
		x = append(x, MasterTask{Task: r.Self()})
	}

	return PipeLine{Tasks: x}
}

func (p PipeLine) Start(ctx context.Context) []TaskResult {

	d := NewDispatcher("Master")
	logrus.Info("Start Simulation PipeLine")
	logrus.Info("Workers: ", config.ParallelConfig.Master)
	logrus.Info("Begin: ", time.Now().Format(time.ANSIC))
	res := d.Dispatch(ctx, config.ParallelConfig.Master, p.Tasks)
	logrus.Info("Finished Simulation PipeLine")
	logrus.Info("End: ", time.Now().Format(time.ANSIC))

	return res
}
