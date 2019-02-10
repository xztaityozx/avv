package cmd

import (
	"context"
	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestDBAccessTask_Run(t *testing.T) {
	as := assert.New(t)
	home, _ := homedir.Dir()
	p := PathJoin(home, "TestDir","db")
	FU.TryMkDir(p)
	p = PathJoin(p, "database.db")

	os.Remove(p)

	r := NewRepositoryFromFile(p)

	dt := DBAccessTask{
		Task:Task{
			Repository:r,
			Failure:10,
			Stage:DBAccess,
			TaskId:9,
			SEED:20,
		},
	}

	res := dt.Run(context.Background())
	as.True(res.Status)
	as.Equal(dt.Task, res.Task)
}
