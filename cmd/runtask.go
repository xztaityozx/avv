package cmd

import (
	"encoding/json"
	"io/ioutil"
	"os"
)


type SimulationResult struct {
	Task   Task
	Status bool
}

type RunTask []Task

// Read Task File
// returns: Task struct, error
func ReadTaskFile(p string) (Task, error) {
	if _, err := os.Stat(p); err != nil {
		return Task{}, err
	}

	b, err := ioutil.ReadFile(p)
	if err != nil {
		return Task{}, err
	}

	var rt Task
	if err := json.Unmarshal(b, &rt); err != nil {
		return Task{}, err
	}

	return rt, nil
}

func GetTasksFromFiles(p ...string) RunTask {
	var rt RunTask

	for _,v := range p {
		if t,err := ReadTaskFile(v); err != nil {
			log.WithError(err).Fatal("Failed Unmarshal Task file: ",v)
		} else {
			rt=append(rt, t)
		}
	}

	return rt
}

func (rt RunTask) BackUp() error {

	m := make(map[string]Repository)
	for _,v := range rt {
		m[v.Repository.Path]=v.Repository
	}

	for _,v:=range m {
		if err := v.DBBackUp();err != nil {
			return err
		}
	}

	return nil
}

func GetTasksFromTaskDir(cnt int) RunTask {
	files, err := ioutil.ReadDir(ReserveDir())
	if err != nil {
		log.WithError(err).Fatal("Can not open: ", ReserveDir())
	}

	var path []string
	for _, v := range files {
		path=append(path, PathJoin(ReserveDir(),v.Name()))
		if cnt > 0 && len(path) < cnt {
			break
		}
	}

	return GetTasksFromFiles(path...)
}