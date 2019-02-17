package cmd

import "time"

type RunSummary struct {
	BeginTime  time.Time
	FinishTime time.Time
	Succeeded  []Task
	Failed     []Task
}
