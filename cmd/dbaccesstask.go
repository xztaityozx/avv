package cmd

import "context"

type DBAccessTask struct {
	Task Task
}

func (da DBAccessTask) Run(parent context.Context) TaskResult {
	ctx, cancel := context.WithCancel(parent)
	defer cancel()

	rr := NewResultRecord(da.Task)
	ch := make(chan error)
	defer close(ch)
	go func() {
		err := rr.Insert(ctx, da.Task.Repository)
		ch<- err
	}()

	select {
	case <-ctx.Done():
	case err := <-ch:
		if err != nil {
			log.WithError(err).WithField("at", "DBAccessTask.Run").Error("Failed Insert to DB")
			return TaskResult{
				Task:   da.Task,
				Status: false,
			}
		}
	}
	return TaskResult{
		Task:da.Task,
		Status:true,
	}
}

func (da DBAccessTask) Self() Task {
	return da.Task
}

func (DBAccessTask) String() string {
	return ""
}
