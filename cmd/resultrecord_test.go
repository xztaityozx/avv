package cmd

import (
	"context"
	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewResultRecord(t *testing.T) {
	as:=assert.New(t)
	actual := NewResultRecord(CountResult{
		Failure:20,
		Task:Task{TaskId:1, SEED:10},
	})

	as.Equal(20, actual.Failure)
	as.Equal(1, actual.TaskId)
	as.Equal(10,actual.Seed)
}

func TestResultRecord_InsertQuery(t *testing.T) {
	as:=assert.New(t)
	expect := "insert into Results(TaskId, Seed, Failure, Date) values (?,?,?,?)"
	actual := ResultRecord{}.InsertQuery()

	as.Equal(expect,actual)
}

func TestResultRecord_SelectQuery(t *testing.T) {
	expect := "select * from Results where TaskId = ?"
	actual := ResultRecord{}.SelectQuery()
	assert.Equal(t,expect,actual)
}

func TestResultRecord_Insert(t *testing.T) {
	as := assert.New(t)
	home, _ := homedir.Dir()
	p := PathJoin(home,"TestDir","db")
	FU.TryMkDir(p)
	p=PathJoin(p,"database.db")

	os.Remove(p)

	r:=NewRepositoryFromFile(p)

	rr := NewResultRecord(CountResult{
		Failure:10,
		Task:Task{
			TaskId:50,
			SEED:20,
			Repository:r,
		}})
	err := rr.Insert(context.Background(),r)
	as.NoError(err)

	t.Run("Check_Select", func(t *testing.T) {
		res, err := rr.Select(context.Background(),r)
		as.NoError(err)
		as.Equal(1,len(res))
		as.Equal(int64(1), res[0].Id)
		as.Equal(int64(50),res[0].TaskId)
		as.Equal(int64(20),res[0].Seed)
		as.Equal(int64(10), res[0].Failure)
	})
}
