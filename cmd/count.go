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
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"syscall"

	pipeline "github.com/mattn/go-pipeline"
	"github.com/spf13/cobra"
)

// countCmd represents the count command
var countCmd = &cobra.Command{
	Use:     "count",
	Aliases: []string{"cnt"},
	Short:   "数え上げします",
	Long: `CSVを受け取って数え上げします

ex)
	avv count --input=file.csv --filter='>=0.4,>=0.4,>=0.4' --out=./out.txt
	
`,
	Run: func(cmd *cobra.Command, args []string) {
		// logger
		l := log.WithField("at", "avv count")

		// 入力ファイル
		inputs, err := cmd.Flags().GetStringSlice("input")
		if err != nil {
			l.Fatal(err)
		}
		// 数え上げクエリ
		query, err := cmd.Flags().GetStringSlice("filter")
		if err != nil {
			l.Fatal(err)
		}

		// クエリが空ならデフォルトを使う
		if len(query) == 0 {
			query = config.Default.PlotPoint.Filters[0].Status
		}

		//Filter struct
		filter := Filter{
			Status: query,
		}

		// 出力先ファイル
		dst, err := cmd.Flags().GetString("out")
		if err != nil {
			l.Fatal(err)
		}

		// 累積和するかどうか
		cum, err := cmd.Flags().GetBool("Cumulative")
		if err != nil {
			l.Fatal(err)
		}

		// ファイル名Validation用正規表現
		reg, err := regexp.Compile("SEED[0-9]+.csv")
		if err != nil {
			l.Fatal(err)
		}

		// 入力ファイルが未指定かallならカレントのファイルを対象にする
		if len(inputs) == 0 || inputs[0] == "all" {
			inputs = []string{}
			wd, _ := os.Getwd()
			fi, err := ioutil.ReadDir(wd)
			if err != nil {
				l.Fatal(err)
			}
			for _, v := range fi {
				inputs = append(inputs, PathJoin(wd, v.Name()))
			}
		}

		var cct []ITask
		// AWKのスクリプト
		script := "BEGIN{sum=0}" + filter.ToAwkStatement(1) + "{sum++}END{print sum}"

		for _, v := range inputs {
			// ファイル名Validation
			base := filepath.Base(v)
			if !reg.MatchString(base) {
				l.Fatal("Does not match file name to `SEED[0-9].csv` regexp")
			}

			// シード値取り出し
			seed, _ := strconv.Atoi(base[5:9])
			cct = append(cct, CommandLineCountTask{SEED: seed, TargetFile: v, ColSize: len(query), Script: script})
		}

		// Dispatcher
		d := NewDispatcher("avv count")

		// キャンセル付きcontext
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// SIGNALをトラップ
		sigCh := make(chan os.Signal)
		defer close(sigCh)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGSTOP)
		go func() {
			<-sigCh
			cancel()
		}()

		// 終了通知チャンネル
		ch := make(chan struct{})
		defer close(ch)

		go func() {
			// Start Task
			res := d.Dispatch(ctx, config.ParallelConfig.CountUp, cct)

			// 結果をSEED値でソート
			sort.Slice(res, func(i, j int) bool {
				return res[i].Task.SEED < res[j].Task.SEED
			})

			// 書き出す
			fp, err := os.OpenFile(dst, os.O_APPEND|os.O_CREATE, 0644)
			defer fp.Close()
			if err != nil {
				l.WithError(err).Error("Failed open out file")
				ch <- struct{}{}
				return
			}

			w := bufio.NewWriter(fp)

			sum := int64(0)
			for _, r := range res {
				if !r.Status {
					l.Error("Find failed task. please retry")
					ch <- struct{}{}
					return
				}

				sum += r.Task.Failure
				w.WriteString(fmt.Sprintf("SEED: %d", r.Task.SEED))
				if cum {
					w.WriteString(fmt.Sprintf("Sum: %d", sum))
				} else {
					w.WriteString(fmt.Sprintf("Failure: %d", r.Task.Failure))
				}
			}
			ch <- struct{}{}
		}()

		select {
		case <-ctx.Done():
		case <-ch:
		}
	},
}

func init() {
	rootCmd.AddCommand(countCmd)

	countCmd.Flags().StringP("out", "o", "result", "specify output file")
	countCmd.Flags().BoolP("Cumulative", "c", false, "累積和で出力します")
	countCmd.Flags().StringSliceP("filter", "f", []string{}, "Query string")
	countCmd.Flags().StringSliceP("input", "i", []string{}, "Input files. if set 'all' or empty, avv pick up all csv file in current directory")
}

type CountTask struct {
	Task Task
}

type CommandLineCountTask struct {
	TargetFile string
	ColSize    int
	SignalName string
	Script     string
	SEED       int
}

// Run for CommandLine Interface Task
func (cct CommandLineCountTask) Run(parent context.Context) TaskResult {
	ctx, cancel := context.WithCancel(parent)
	defer cancel()

	ch := make(chan int64)
	defer close(ch)
	ech := make(chan error)
	defer close(ech)

	go func() {
		p, err := ShapingCSV(cct.TargetFile, cct.SignalName, cct.ColSize)
		if err != nil {
			ech <- err
			return
		}
		out, err := exec.Command("awk", cct.Script, p).Output()
		if err != nil {
			ech <- err
			return
		}

		f, err := strconv.Atoi(string(out))
		if err != nil {
			ech <- err
			return
		}

		ch <- int64(f)
	}()

	l := log.WithField("at", "CommandLineCountTask.Run")
	select {
	case <-ctx.Done():
		l.Warn("Canceled By Context")
		return TaskResult{Status: false}
	case err := <-ech:
		l.Error(err)
		return TaskResult{Status: false}
	case f := <-ch:
		return TaskResult{
			Task: Task{
				SEED:    cct.SEED,
				Failure: int64(f),
			},
			Status: true,
		}
	}

}

func (cct CommandLineCountTask) Self() Task {
	return Task{}
}

func (CommandLineCountTask) String() string {
	return ""
}

func (ct CountTask) Run(parent context.Context) TaskResult {
	ctx, cancel := context.WithCancel(parent)
	defer cancel()

	ch := make(chan int64)
	defer close(ch)
	ech := make(chan error)
	defer close(ech)

	// fork Task
	go func() {
		f, err := ct.CountUp()
		if err != nil {
			ech <- err
		} else {
			ch <- f
		}
	}()

	// return Result
	select {
	case <-ctx.Done():
		return TaskResult{
			Status: false,
			Task:   ct.Task,
		}
	case err := <-ech:
		log.WithField("at", "CountTask.Run").WithError(err).Error("Failed CountUp")
		return TaskResult{
			Status: false,
			Task:   ct.Task,
		}
	case rec := <-ch:
		ct.Task.Failure = rec
		return TaskResult{
			Status: true,
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
func (ct CountTask) CountUp() (failure int64, err error) {
	resultDir := ct.Task.SimulationDirectories.ResultDir
	// Can not find resultDir
	if _, err := os.Stat(resultDir); err != nil {
		return -1, err
	}

	// find target files
	var targets []string
	for _, v := range ct.Task.PlotPoint.Filters {
		p := PathJoin(resultDir, v.SignalName, fmt.Sprintf("SEED%05d.csv", ct.Task.SEED))
		if _, err := os.Stat(p); err != nil {
			return -1, err
		}

		// Convert file format
		out, err := ShapingCSV(p, v.SignalName, len(v.Status))
		if err != nil {
			return -1, err
		}

		targets = append(targets, out)
	}

	out, err := pipeline.Output(
		[]string{"paste", strings.Join(targets, " ")},
		[]string{"awk", ct.Task.PlotPoint.GetAwkScript()})

	if err != nil {
		log.WithField("at", "CountTask.CountUp").Error(string(out))
		return -1, err
	}

	failure, err = strconv.ParseInt(string(out), 10, 64)
	return
}

// Convert CSV file that generated by WaveView to one tran's record per one line
//returns: output file path, error
func ShapingCSV(p, signalName string, n int) (string, error) {
	out, err := ioutil.ReadFile(p)
	if err != nil {
		return "", err
	}

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
