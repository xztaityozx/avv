// Copyright © 2019 xztaityozx
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"os/exec"
)

type SimulationTask struct {
	Task Task
}

func (s SimulationTask) Run(parent context.Context) TaskResult {
	t := s.Task

	// Make Simulation Parameter and Directories
	t.MkDir()
	t.SimulationFiles.AddFile.Make(t.SimulationDirectories.BaseDir)
	t.MakeSPIScript()

	cmdStr := t.GetSimulationCommand()
	log.WithField("at", "Task.RunSimulation").Info("Command: ", cmdStr)

	command := exec.Command("bash", "-c", cmdStr)

	// Run Simulation

	l := log.WithField("at","SimulationTask").WithField("seed",fmt.Sprint(t.SEED))
	l.Info("Start Simulation")

	out, err := command.Output()
	if err != nil {
		logfile := PathJoin(t.SimulationDirectories.DstDir, "hspice.log")
		log.WithField("at", "SimulationTask.Run()").
			WithError(err).
			Error("Failed Simulation:" + string(out) + " hspice.log=" + FU.Cat(logfile))

		return TaskResult{
			Status:false,
			Task:t,
		}
	}

	l.Info("Simulation finished")

	l.Info("Start File Check")
	files, err := ioutil.ReadDir(t.SimulationDirectories.DstDir)
	if err != nil {
		l.Error(err)
		return TaskResult{
			Status:false,
			Task:t,
		}
	}
	if len(files) < t.Times {
		l.Error("波形データが少なすぎます。")
		return TaskResult{
			Status:false,
			Task:t,
		}
	}
	l.Info("file check finished")

	return TaskResult{
		Task:   t,
		Status: err == nil,
	}
}

func (s SimulationTask) String() string {
	return ""
}

func (s SimulationTask) Self() Task {
	return s.Task
}

func (hc HSPICEConfig) GetCommand(spi string) string {
	return fmt.Sprintf("%s %s -i %s -o ./hspice &> ./hspice.log", hc.Command, hc.Option, spi)
}
