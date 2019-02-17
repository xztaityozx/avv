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

	task := GetTasksFromTaskDir(10)

	as.Equal(10, len(task))

	for _, v := range task {
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
	rt := GetTasksFromFiles(path, path, path)

	for _, v := range rt {
		as.Equal(config.Default, v)
	}
}
