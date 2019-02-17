package cmd

import (
	"context"
	"fmt"
)

type Transistor struct {
	TransistorId int64
	Threshold    float64
	Sigma        float64
	Deviation    float64
}

// Compare func for Transistor struct
func (t Transistor) Compare(s Transistor) bool {
	return t.Sigma == s.Sigma && t.Deviation == s.Deviation && t.Threshold == s.Threshold && t.TransistorId == s.TransistorId
}

func (t Transistor) StringPrefix(prefix string) string {
	return fmt.Sprintf("%s%.4f-Sigma%.4f", prefix, t.Threshold, t.Sigma)
}
func (t Transistor) SelectQuery() string {
	return fmt.Sprintf("select * from Transistor where Deviation = ? and Threshold = ? and Sigma = ?")
}

func (t Transistor) Select(ctx context.Context, repository Repository) (Transistor, error) {
	db, err := repository.Connect()
	defer db.Db.Close()
	if err != nil {
		return Transistor{}, err
	}

	var rt Transistor
	err = db.WithContext(ctx).SelectOne(&rt, t.SelectQuery(), t.Deviation, t.Threshold, t.Sigma)

	return rt, err
}

func (t Transistor) InsertQuery() string {
	return "insert or ignore into Transistor(Deviation, Threshold, Sigma) values (?,?,?)"
}

func (t Transistor) Insert(ctx context.Context, repository Repository) error {
	db, err := repository.Connect()
	defer db.Db.Close()
	if err != nil {
		return err
	}

	_, err = db.WithContext(ctx).Exec(t.InsertQuery(), t.Deviation, t.Threshold, t.Sigma)
	return err
}
