package cmd

import (
	"context"
	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestTaskGroup_InsertQuery(t *testing.T) {
	as := assert.New(t)
	expect := "insert or ignore into Groups(ParamsId, SeedStart, SeedEnd, Date) values (?,?,?,?)"
	actual := TaskGroup{}.InsertQuery()
	as.Equal(expect, actual)

}

func TestTaskGroup_SelectQuery(t *testing.T) {
	as := assert.New(t)
	expect := "select * from Groups where ParamsId = ? and SeedStart = ? and SeedEnd = ? and Date = ?"
	actual := TaskGroup{}.SelectQuery()

	as.Equal(expect, actual)
}

func TestTaskGroup_Insert(t *testing.T) {
	as := assert.New(t)
	home, _ := homedir.Dir()
	p := PathJoin(home, "TestDir", "db")
	FU.TryMkDir(p)
	p = PathJoin(p, "database.db")

	os.Remove(p)
	r := NewRepositoryFromFile(p)

	tg := TaskGroup{
		Date:      time.Now(),
		SeedEnd:   2000,
		SeedStart: 1,
		ParamsId:  1,
	}

	err := tg.Insert(context.Background(), r)
	as.NoError(err)
	t.Run("Check_Select", func(t *testing.T) {
		rec, err := tg.Select(context.Background(), r)
		tg.TaskId = 1
		as.NoError(err)
		as.True(tg.Compare(rec))
	})
}

func TestMakeRequest_NewTaskGroup(t *testing.T) {
	as := assert.New(t)
	home, _ := homedir.Dir()
	p := PathJoin(home, "TestDir", "db")
	FU.TryMkDir(p)
	p = PathJoin(p, "database.db")

	os.Remove(p)
	r := NewRepositoryFromFile(p)

	mt := MakeRequest{
		SEED: SEED{
			Start: 1,
			End:   2000,
		},
		Task: Task{Repository: r},
	}

	tg, err := mt.NewTaskGroup(context.Background())
	as.NoError(err)
	as.Equal(int64(1), tg.TaskId)

}
