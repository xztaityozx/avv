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
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	wvparser "github.com/xztaityozx/go-wvparser"
	"golang.org/x/xerrors"
)

// countCmd represents the count command
var countCmd = &cobra.Command{
	Use:     "count",
	Aliases: []string{"cnt"},
	Short:   "数え上げします",
	Long: `CSVを受け取って数え上げします


`,
	Run: func(cmd *cobra.Command, args []string) {
		ofile, _ := cmd.Flags().GetString("out")
		filter, _ := cmd.Flags().GetStringSlice("filter")
		parallel, _ := rootCmd.PersistentFlags().GetInt("Parallel")

		var task []ITask
		for _, v := range args {
			task = append(task, NewCountCLITask(v, filter))
		}

		d := NewDispatcher("Counter")
		res := d.Dispatch(context.Background(), parallel, task)

		var box = map[string]int64{}

		for _, v := range res {
			box[v.Task.SimulationFiles.Self] = v.Task.Failure
		}

		if len(ofile) == 0 {
			for i, v := range box {
				fmt.Println(i, ": ", v)
			}
		} else {
			str := ""
			for i, v := range box {
				str = fmt.Sprintf("%s%s : %d\n", str, i, v)
			}
			ioutil.WriteFile(ofile, []byte(str), 0644)
		}

	},
}

func init() {
	rootCmd.AddCommand(countCmd)
	countCmd.Flags().StringSlice("filter", []string{}, "フィルター")
	countCmd.Flags().StringP("out", "o", "", "出力ファイル")
}

type CountTask struct {
	Task Task
}

func (ct CountTask) Run(parent context.Context) TaskResult {
	ctx, cancel := context.WithCancel(parent)
	defer cancel()

	ech := make(chan error)
	defer close(ech)

	// fork Task
	go func() {
		ech <- ct.CountUp()
	}()

	// return Result
	select {
	case <-ctx.Done():
		return TaskResult{
			Status: false,
			Task:   ct.Task,
		}
	case err := <-ech:
		if err != nil {
			log.WithField("at", "CountTask.Run").WithError(err).Error("Failed CountUp")
		}
		return TaskResult{
			Status: err == nil,
			Task:   ct.Task,
		}
	}
}

func (ct CountTask) String() string {
	return ""
}

func (ct CountTask) Self() Task {
	return ct.Task
}

// CountUp Aggregate failure from csv file which generated by WaveView
// returns: number of failure, error
func (ct CountTask) CountUp() error {
	resultDir := ct.Task.SimulationDirectories.ResultDir
	// Can not find resultDir
	if _, err := os.Stat(resultDir); err != nil {
		return err
	}

	for _, v := range ct.Task.PlotPoint.Signals {
		p := PathJoin(resultDir, v, fmt.Sprintf("SEED%05d.csv", ct.Task.SEED))
		if _, err := os.Stat(p); err != nil {
			return err
		}

		csv, err := wvparser.WVParser{FilePath: p}.Parse()
		if err != nil {
			return xerrors.Errorf("Failed Parse: %w", err)
		}

		if err := config.Server.Insert(ct.Task.Vtn, ct.Task.Vtp, ct.Task.SEED, &csv); err != nil {
			return err
		}
	}
	return nil
}

// Convert CSV file that generated by WaveView to one tran's record per one line
//returns: output file path, error
func ShapingCSV(p, signalName string, n int) (string, error) {
	out, err := ioutil.ReadFile(p)
	if err != nil {
		return "", err
	}

	tmpavv := PathJoin("/tmp/avv")
	FU.TryMkDir(tmpavv)
	tmp, err := ioutil.TempFile("/tmp/avv", signalName+".*.csv")
	if err != nil {
		return "", err
	}
	w := bufio.NewWriter(tmp)
	defer w.Flush()

	idx := 0
	for _, line := range bytes.Split(out, []byte("\n")) {
		if len(line) == 0 || line[0] == byte('#') || line[0] == byte('T') {
			continue
		}

		data := bytes.Split(bytes.Replace(line, []byte(" "), []byte(""), -1), []byte(","))[1]
		_, err := w.Write(data)
		if err != nil {
			return "", err
		}
		idx++
		if idx%n == 0 {
			_, err := w.WriteString("\n")
			if err != nil {
				return "", err
			}
		} else {
			_, err := w.WriteString(" ")
			if err != nil {
				return "", err
			}
		}
	}

	return tmp.Name(), nil
}
