package cmd

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

type RunTask struct {
	Tasks      []Task
	RunSummary RunSummary
}

type SimulationResult struct{
	Task Task
	Status bool
}

func (rt RunTask) Run() {
	l := log.WithField("at","RunTask.Run")


	simWorker := func(tasks []Task) <- chan SimulationResult {
		l.Info("Start HSPICE Simulations")
		rec := make(chan SimulationResult, config.ParallelConfig.HSPICE)
		for _, v := range tasks {
			go func(t Task) {
				if err := t.RunSimulation(); err != nil {
					rec <- SimulationResult{
						Task:t,
						Status:false,
					}
				} else {
					rec <- SimulationResult{
						Task:t,
						Status:true,
					}
				}
			}(v)
		}

		return rec
	}

	extWorker := func(tasks []Task) <- chan SimulationResult {
		l.Info("Start Extract")
		rec := make(chan SimulationResult, config.ParallelConfig.WaveView)
		for _, v := range tasks {
			go func(t Task) {
				if err := t.RunExtract(); err != nil {
					rec <- SimulationResult{
						Task:t,
						Status:false,
					}
				} else {
					rec <- SimulationResult{
						Task:t,
						Status:true,
					}
				}
			}(v)
		}
		return rec
	}
	
}

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

// Get Tasks
// returns: error, []Task(ref)
func (rt *RunTask) GetTasks(cnt int) error {
	taskDir := ReserveDir()
	files, err := ioutil.ReadDir(taskDir)
	if err != nil {
		return err
	}

	// Min(cnt,len(files))
	if cnt > len(files) || cnt < 0 {
		cnt = len(files)
	}

	for i := 0; i < cnt; i++ {
		v := files[i]
		fulPath := PathJoin(taskDir, v.Name())
		if t, err := ReadTaskFile(fulPath); err != nil {
			log.WithField("at", "RunTask.GetTasks").Error("Failed read task file: ", v.Name(), "This file will move to ", DustDir())
			if err := os.Rename(fulPath, PathJoin(DustDir(), v.Name())); err != nil {
				return err
			}
		} else {
			rt.Tasks = append(rt.Tasks, t)
		}
	}

	return nil
}

// Get Tasks
// returns: error, []Task(ref)
func (rt *RunTask) GetTaskFromFiles(f ...string) error {

	for _, v := range f {
		if t, err := ReadTaskFile(v); err != nil {
			log.WithField("at", "RunTask.GetRunTaskFromFiles").Error("Failed read task file: ", v, "This file will move to", DustDir())
			if err := os.Rename(v, PathJoin(DustDir(), filepath.Base(v))); err != nil {
				return err
			}
		} else {
			rt.Tasks = append(rt.Tasks, t)
		}
	}
	return nil
}
