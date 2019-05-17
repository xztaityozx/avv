package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type ExtractTask struct {
	Task Task
}

// Generate wv Command
func (wvc WaveViewConfig) GetCommand(ace string) string {
	return fmt.Sprintf("%s -k -ace_no_gui %s &> wv.log", wvc.Command, ace)
}

// Generate extract command
func (t ExtractTask) GetExtractCommand() string {
	var rt []string

	// append cd command
	rt = append(rt, t.Task.GetCdCommand())
	// append waveview command
	rt = append(rt, config.WaveView.GetCommand(t.Task.SimulationFiles.ACEScript))

	return strings.Join(rt, " ")
}

func (t ExtractTask) String() string {
	return fmt.Sprint()
}

// Run extract command with WaveView
// returns: errors
func (t ExtractTask) Run(parent context.Context) TaskResult {

	l := log.WithField("at", "ExtractTask")

	// make ace script
	ace, err := t.Task.PlotPoint.MkACEScript(t.Task.SimulationDirectories.DstDir)
	l.Info(t.Task.SimulationDirectories.DstDir)
	if err != nil {
		l.WithError(err).Error("Failed make ACEScript")
		return TaskResult{
			Task:   t.Task,
			Status: false,
		}
	}
	t.Task.SimulationFiles.ACEScript = ace

	// Make results.xml
	if path, err := t.Task.MakeResultsXml(); err != nil {
		l.WithError(err).Error("Failed make results.xml")
		return TaskResult{
			Task:   t.Task,
			Status: false,
		}
	} else {
		t.Task.SimulationFiles.ResultsXML = path
	}

	// Make resultsMap.xml
	if path, err := t.Task.MakeMapXml(); err != nil {
		l.WithError(err).Error("Failed make resultsMap.xml")
		return TaskResult{
			Task:   t.Task,
			Status: false,
		}
	} else {
		t.Task.SimulationFiles.ResultsMapXML = path
	}

	cmdStr := t.GetExtractCommand()
	command := exec.Command("bash", "-c", cmdStr)
	_, err = command.CombinedOutput()
	if err != nil {
		logfile := PathJoin(t.Task.SimulationDirectories.DstDir, "wv.log")
		l.WithError(err).Error("Failed Extract: " + FU.Cat(logfile))
		return TaskResult{
			Task:   t.Task,
			Status: false,
		}
	}

	for _, v := range t.Task.PlotPoint.Signals {
		// csvをそれぞれの信号線のディレクトリに移動させる
		oldPath := PathJoin(t.Task.SimulationDirectories.DstDir, v+".csv")
		FU.TryMkDir(PathJoin(t.Task.SimulationDirectories.ResultDir, v))
		newPath := PathJoin(t.Task.SimulationDirectories.ResultDir, v, fmt.Sprintf("SEED%05d.csv", t.Task.SEED))

		if err := os.Rename(oldPath, newPath); err != nil {
			l.WithError(err).Error("Failed Rename ", oldPath, " to ", newPath)
			return TaskResult{
				Task:   t.Task,
				Status: false,
			}
		}

		t.Task.ResultCSV = append(t.Task.ResultCSV, newPath)
	}

	return TaskResult{
		Task:   t.Task,
		Status: true,
	}
}

func (t ExtractTask) Self() Task {
	return t.Task
}
