package cmd

import (
	"context"
	"os"

	"github.com/xztaityozx/go-wvparser"
)

type CountCLITask struct {
	Path    string
	Counter wvparser.Counter
}

func NewCountCLITask(path string, f []string) CountCLITask {
	return CountCLITask{
		Path:    path,
		Counter: wvparser.NewCounter(f...),
	}
}

func (cct CountCLITask) Run(context context.Context) TaskResult {

	ch := make(chan int64)
	defer close(ch)

	go func() {
		path := cct.Path
		if _, err := os.Stat(path); err != nil {
			log.WithError(err).Fatal("Failed count sub command")
		}

		csv, err := wvparser.WVParser{FilePath: path}.Parse()
		if err != nil {
			log.WithError(err).Fatal("Failed Parse")
		}

		ch <- cct.Counter.Aggregate(csv)
	}()

	select {
	case <-context.Done():
		return TaskResult{Status: false}
	case res := <-ch:
		return TaskResult{
			Task: Task{
				SimulationFiles: SimulationFiles{
					Self: cct.Path,
				},
				Failure: res,
			},
			Status: true,
		}
	}

}

func (cct CountCLITask) Self() Task {
	return Task{}
}
