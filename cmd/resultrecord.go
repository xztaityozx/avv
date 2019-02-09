package cmd

import (
	"time"
)

type (
	ResultRecord struct {
		Id      int
		TaskId  int64
		Seed    int64
		Failure int64
		Date    time.Time
	}
)

// NewResultRecord create ResultRecord struct
// args: succeeded TaskResult
// returns: ResultRecord struct not inserted yet
func NewResultRecord(result CountResult) ResultRecord {
	return ResultRecord{
		TaskId:result.Task.TaskId,
		Date:time.Now(),
		Seed:int64(result.Task.SEED),
		Failure:result.Failure,
	}
}
