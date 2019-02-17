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

	t.Run("Dispatch", func(t *testing.T) {
		d := NewDispatcher("DBAccess")
		var dts []ITask
		for i:=0;i<20;i++ {
			dts=append(dts, dt)
		}

		res := d.Dispatch(context.Background(), 4, dts)

		as.Equal(20, d.Size)
		as.Equal(20,len(res))

		for _, v := range res{
			as.Equal(TaskResult{Task:dt.Task, Status:true}, v)
		}
		
	})
}

