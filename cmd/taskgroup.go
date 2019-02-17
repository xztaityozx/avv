package cmd

import (
	"context"
	"time"
)

type TaskGroup struct {
	TaskId    int64
	ParamsId  int64
	SeedStart int64
	SeedEnd   int64
	Date      time.Time
}

func (m MakeRequest) NewTaskGroup(ctx context.Context) (TaskGroup, error) {
	repository := m.Task.Repository
	p,err := m.Task.NewParameter(ctx, repository)
	if err != nil {
		return TaskGroup{},err
	}

	err = p.Insert(ctx,repository)
	if err != nil {
		return TaskGroup{}, err
	}

	p, err = p.Select(ctx, repository)
	if err != nil {
		return TaskGroup{}, err
	}

	tg := TaskGroup{
		SeedEnd:int64(m.SEED.End),
		SeedStart:int64(m.SEED.Start),
		ParamsId:p.ParamsId,
		Date:time.Now(),
	}
	err = tg.Insert(ctx, repository)
	if err != nil {
		return TaskGroup{}, err
	}

	return tg.Select(ctx,repository)
}

func (t TaskGroup) Compare(s TaskGroup) bool {
	return t.TaskId == s.TaskId && t.ParamsId == s.ParamsId && t.SeedStart == s.SeedStart && t.SeedEnd == s.SeedEnd
}

func (TaskGroup) InsertQuery() string {
	return "insert or ignore into Groups(ParamsId, SeedStart, SeedEnd, Date) values (?,?,?,?)"
}

func (t TaskGroup) Insert(ctx context.Context, repository Repository) error {
	db, err := repository.Connect()
	defer db.Db.Close()
	if err != nil {
		return err
	}

	_, err = db.WithContext(ctx).Exec(t.InsertQuery(), t.ParamsId, t.SeedStart, t.SeedEnd, t.Date)
	if err != nil {
		return err
	}

	return nil
}

func (TaskGroup) SelectQuery() string {
	return "select * from Groups where ParamsId = ? and SeedStart = ? and SeedEnd = ? and Date = ?"
}

func (t TaskGroup) Select(ctx context.Context, repository Repository) (TaskGroup, error) {
	var rt TaskGroup

	db, err := repository.Connect()
	defer db.Db.Close()
	if err != nil {
		return rt, err
	}

	err = db.WithContext(ctx).SelectOne(&rt, t.SelectQuery(), t.ParamsId, t.SeedStart, t.SeedEnd, t.Date)

	return rt, nil
}
