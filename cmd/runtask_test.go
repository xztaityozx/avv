package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestReadTaskFile(t *testing.T) {
	as := assert.New(t)
	home, _ := homedir.Dir()
	testdir := PathJoin(home, "TestDir")
	FU.TryMkDir(testdir)

	path := PathJoin(testdir, "test.json")

	if b, err := json.Marshal(config.Default); err != nil {
		as.FailNow("", err)
	} else {
		ioutil.WriteFile(path, b, 0644)
	}
	var actual Task
	if b, err := ioutil.ReadFile(path); err != nil {
		as.FailNow("", err)
	} else {
		err := json.Unmarshal(b, &actual)
		if err != nil {
			as.FailNow("Unmarshal error: ", err)
		}
	}

	as.Equal(config.Default, actual)

	ioutil.WriteFile(path, []byte("abc"), 0644)

	if _, err := ReadTaskFile(path); err == nil {
		as.FailNow("Invalid Marshal")
	}

	os.Remove(path)
}

func TestRunTask_GetTasks(t *testing.T) {
	var rt RunTask
	home, _ := homedir.Dir()
	config.TaskDir = PathJoin(home, "TestDir", "Task")
	FU.TryMkDir(config.TaskDir)
	FU.TryMkDir(ReserveDir())

	b, _ := json.Marshal(config.Default)

	for i := 0; i < 10; i++ {
		p := PathJoin(ReserveDir(), fmt.Sprintf("%d.json", i))
		ioutil.WriteFile(p, b, 0644)
	}
	as := assert.New(t)

	err := rt.GetTasks(10)
	if err != nil {
		as.FailNow("failed GetTasks: ", err)
	}

	as.Equal(10, len(rt.Tasks))

	for _, v := range rt.Tasks {
		as.Equal(config.Default, v)
	}
}

func TestRunTask_GetTaskFromFiles(t *testing.T) {
	home, _ := homedir.Dir()
	config.TaskDir = PathJoin(home, "TestDir", "Task")
	FU.TryMkDir(config.TaskDir)
	FU.TryMkDir(ReserveDir())

	b, _ := json.Marshal(config.Default)

	path := PathJoin(config.TaskDir, "test.json")
	ioutil.WriteFile(path, b, 0644)

	as := assert.New(t)
	var rt RunTask
	err := rt.GetTaskFromFiles(path, path, path)
	if err != nil {
		as.FailNow("RunTask.GetTaskFromFiles", err)
	}

	for _, v := range rt.Tasks {
		as.Equal(config.Default, v)
	}
}
