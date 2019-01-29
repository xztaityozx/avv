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
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

// simulateCmd represents the simulate command
var simulateCmd = &cobra.Command{
	Use:   "simulate",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("simulate called")
	},
}

func init() {
	rootCmd.AddCommand(simulateCmd)
}

type SimulateTask struct {
	Task Task
}

func NewSimulateTask() SimulateTask {
	return SimulateTask{
		Task: config.Default,
	}
}

func (st SimulateTask) Run() error {
	task := st.Task

	// Make Simulation Parameter and Directories
	task.MkDir()
	task.SimulationFiles.AddFile.Make(task.SimulationDirectories.BaseDir)
	task.MakeSPIScript()

	command := st.GetCommand()
	log.WithField("at", "SimulateTask.Run").Info(command)

	c := exec.Command("bash", "-c", command)

	// Run Simulation
	out, err := c.CombinedOutput()
	if err != nil {
		return errors.New("Failed Simulation\n")
	} else {
		log.WithField("at", "SimulateTask.Run").Info(string(out))
	}

	return nil
}

func (hc HSPICEConfig) GetCommand(spi string) string {
	return fmt.Sprintf("%s %s -i %s -o ./hspice &> ./hspice.log", hc.Command, hc.Option, spi)
}

func (st SimulateTask) GetCommand() string {
	var rt []string

	// append cd command
	rt = append(rt, fmt.Sprintf("cd %s &&", st.Task.SimulationDirectories.DstDir))
	// append hspice command
	rt = append(rt, config.HSPICE.GetCommand(st.Task.SimulationFiles.SPIScript))
	return strings.Join(rt, " ")
}
