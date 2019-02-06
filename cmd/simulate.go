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
	"github.com/spf13/cobra"
	"os/exec"
)

// simulateCmd represents the simulate command
var simulateCmd = &cobra.Command{
	Use:   "simulate",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		//if len(args[0]) == 0 {
		//	log.WithField("command", "simulate").Fatal("Invalid file path")
		//}
		//
		//task, err := ReadTaskFile(args[0])
		//if err != nil {
		//	log.WithField("command", "simulate").Fatal(err)
		//}
		//
		//log.WithField("command", "simulate").Info("Start Simulation at ", time.Now().Format(time.RFC3339))
		//s := spinner.New(spinner.CharSets[36], time.Millisecond*500)
		//s.FinalMSG = "Finished!!"
		//s.Suffix = "Running..."
		//s.Start()
		//err = SimulationTask{Task:task}.Run()
		//if err != nil {
		//	log.WithField("command", "simulate").Fatal(err)
		//}
		//s.Stop()

	},
}

func init() {
	rootCmd.AddCommand(simulateCmd)
}

type SimulationTask struct {
	Task Task
}

func (s SimulationTask) Run(parent context.Context) Result {
	t := s.Task

	// Make Simulation Parameter and Directories
	t.MkDir()
	t.SimulationFiles.AddFile.Make(t.SimulationDirectories.BaseDir)
	t.MakeSPIScript()

	cmdStr := t.GetSimulationCommand()
	log.WithField("at", "Task.RunSimulation").Info("Command: ", cmdStr)

	command := exec.Command("bash", "-c", cmdStr)

	// Run Simulation
	out, err := command.CombinedOutput()
	if err != nil {
		logfile := PathJoin(t.SimulationDirectories.DstDir, "hspice.log")
		log.WithField("at", "SimulationTask.Run()").
			WithError(err).
			Error("Failed Simulation:" + string(out) + " hspice.log=" + FU.Cat(logfile))
	}

	return Result{
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
