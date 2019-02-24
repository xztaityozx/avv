package cmd

import "context"

type MasterTask struct {
	Task   Task
	DBPipe chan DBAccessTask
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
			t.Stage = DBAccess
			return DBAccessTask{Task: t}
		} else {
			t.Stage = Remove
			return RemoveTask{Task: t}
		}
	}

	for t.Self().Stage != DBAccess {
		res := t.Run(ctx)
		if !res.Status {
			return res
		}

		t = next(res.Task)
	}

	mt.DBPipe <- DBAccessTask{Task: t.Self()}

	res := RemoveTask{Task: t.Self()}.Run(ctx)

	return res
}
