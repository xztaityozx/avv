package cmd

import "context"

type MasterTask struct {
	Task   Task
}

func (mt MasterTask) Self() Task {
	return mt.Task
}

func (mt MasterTask) String() string {
	return ""
}

func (mt MasterTask) Run(ctx context.Context) TaskResult {

	t := mt.Task.GetWrapper()
	next := func(t Task) ITask {
		if t.Stage == HSPICE {
			t.Stage = WaveView
			return ExtractTask{Task: t}
		} else if t.Stage == WaveView {
			t.Stage = CountUp
			return CountTask{Task: t}
		} else if t.Stage == CountUp {
			t.Stage = Remove
			return RemoveTask{Task: t}
		} else {
			t.Stage = DBAccess
			return DBAccessTask{Task: t}
		}
	}

	retry := map[string]int {
		string(HSPICE): config.RetryConfig.HSPICE,
		string(WaveView): config.RetryConfig.WaveView,
		string(CountUp): config.RetryConfig.CountUp,
		string(Remove): 0,
	}

	for t.Self().Stage != DBAccess {
		var res TaskResult
		var cnt int
		for ok := true; ok; ok = !res.Status && cnt < retry[string(t.Self().Stage)] {
			res = t.Run(ctx)
			if !res.Status {
				log.WithField("at","MasterTask.Run").Warn("Failed Task Step at: ", t.Self().Stage, " Retry...(",cnt,")")
				cnt++
			}
		}

		if !res.Status {
			return TaskResult{
				Task:t.Self(),
				Status:false,
			}
		}

		t = next(res.Task)
	}

	return TaskResult{
		Task:t.Self(),
		Status:true,
	}
}
