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
	"github.com/spf13/cobra"
	"github.com/xztaityozx/avv/extract"
	"github.com/xztaityozx/avv/pipeline"
	"github.com/xztaityozx/avv/push"
	"github.com/xztaityozx/avv/simulation"
	"github.com/xztaityozx/avv/task"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "シミュレーションを実行します",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		n, _ := cmd.Flags().GetInt("count")
		all, _ := cmd.Flags().GetBool("all")
		skip, _ := cmd.Flags().GetBool("skip")

		x, _ := cmd.Flags().GetInt("simulateParallel")
		y, _ := cmd.Flags().GetInt("extractParallel")

		taskdir := config.Default.TaskDir()

		files, err := ioutil.ReadDir(taskdir)
		if err != nil {
			log.Fatal(err)
		}

		size := len(files)
		if size > n {
			size = n
		}

		p := pipeline.New(size)
		source := make(chan task.Task, p.Total)
		for i := 0; i < len(files) && (i < n || all); i++ {
			t, err := task.Unmarshal(filepath.Join(taskdir, files[i].Name()))
			if err != nil {
				log.Fatal(err)
			}

			source <- t
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// set cancel signals
		sigCh := make(chan os.Signal)
		defer close(sigCh)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGSTOP)
		go func() {
			<-sigCh
			cancel()
		}()

		first := p.AddStage(x, source, "simulation", func(ctx context.Context, t task.Task) (i task.Task, e error) {
			h := simulation.HSPICE{
				Path:    config.HSPICE.Path,
				Options: config.HSPICE.Options,
			}
			e = h.Invoke(ctx, t)
			return t, e
		})

		second := p.AddStage(y, first, "extract", func(ctx context.Context, t task.Task) (i task.Task, e error) {
			w := extract.WaveView{
				Path: config.WaveView.Path,
			}
			e = w.Invoke(ctx, t)
			return t, e
		})

		// error channel
		errCh := make(chan error)
		defer close(errCh)

		// Start PipeLine
		go func() {
			errCh <- p.Start(ctx)
		}()

		select {
		case <-ctx.Done():
			log.Fatal("canceled")
		case err := <-errCh:
			if err != nil {
				log.Fatal(err)
			}
		}

		if skip {
			log.Info("Skip push")
			return
		}

		// push
		results := map[push.TaaResultKey][]string{}
		for v := range second {
			key := push.TaaResultKey{
				Vtn:    v.Parameters.Vtn,
				Vtp:    v.Parameters.Vtp,
				Sweeps: v.Sweeps,
			}
			if results[key] == nil {
				results[key] = []string{}
			}

			results[key] = append(results[key], v.Files.ResultFile)
		}

		taa := push.Taa{
			TaaPath:    config.Taa.Path,
			ConfigFile: config.Taa.ConfigFile,
			Parallel:   config.Taa.Parallel,
		}

		log.Info("Start pushing")
		for key, val := range results {
			err := taa.Invoke(ctx, key.Vtn, key.Vtp, key.Sweeps, val)
			if err != nil {
				log.Fatal(err)
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().IntP("count", "n", 1, "number of task")
	runCmd.Flags().Bool("all", false, "")

	runCmd.Flags().IntP("simulateParallel", "x", 1, "HSPICEの並列数です")
	runCmd.Flags().IntP("extractParallel", "y", 1, "WaveViewの並列数です")

	runCmd.Flags().Bool("skip", false, "pushをスキップします")
}
