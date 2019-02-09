package cmd

import (
	"context"
	"time"
)

type (
	ResultRecord struct {
		Id      int64
		TaskId  int64
		Seed    int64
		Failure int64
		Date    time.Time
	}
)

func (ResultRecord) InsertQuery() string {
	return "insert into Results(TaskId, Seed, Failure, Date) values (?,?,?,?)"
}

func (rr ResultRecord) Select(ctx context.Context, r Repository) ([]ResultRecord, error) {
	var rt []ResultRecord

	db, err := r.Connect()
	defer db.Db.Close()
	if err != nil {
		return nil, err
	}

	_, err = db.WithContext(ctx).Select(&rt,rr.SelectQuery(),rr.TaskId)

	return rt, nil
}

func (rr ResultRecord) Insert(ctx context.Context, repository Repository) error {
	db, err :=repository.Connect()
	defer db.Db.Close()
	if err != nil {
		return err
	}

	_, err = db.WithContext(ctx).Exec(rr.InsertQuery(), rr.TaskId,rr.Seed,rr.Failure,rr.Date)
	return err
}

func (ResultRecord) SelectQuery() string {
	return "select * from Results where TaskId = ?"
}

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

