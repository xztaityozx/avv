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
	"github.com/briandowns/spinner"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/xztaityozx/avv/extract"
	"github.com/xztaityozx/avv/push"
	"github.com/xztaityozx/avv/remove"
	"github.com/xztaityozx/avv/simulation"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/xztaityozx/avv/pipeline"
	"github.com/xztaityozx/avv/task"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "シミュレーションを実行します",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		n, _ := cmd.Flags().GetInt("count")
		all, _ := cmd.Flags().GetBool("all")

		x, _ := cmd.Flags().GetInt("simulateParallel")
		y, _ := cmd.Flags().GetInt("extractParallel")
		z, _ := cmd.Flags().GetInt("pushParallel")
		slack, _ := cmd.Flags().GetBool("slack")
		keepcsv, _ := cmd.Flags().GetBool("keepcsv")

		taskdir := config.Default.TaskDir()

		files, err := ioutil.ReadDir(taskdir)
		if err != nil {
			log.WithError(err).Fatal("Failed read dir")
		}

		size := len(files)
		if size > n && !all {
			size = n
		}

		// Find task files
		var box []task.Task
		for i := 0; i < size; i++ {
			// /path/to/taskFile
			path := filepath.Join(taskdir, files[i].Name())

			log.Info("Unmarshal :", path)

			t, err := task.Unmarshal(path)
			if err != nil {
				log.WithError(err).Fatal("Failed unmarshal task file")
			}
			t.Files.TaskFile = path

			box = append(box, t)
		}

		// generate pipeline struct
		p := pipeline.New(len(box), config.MaxRetry)

		// push task to source chan
		source := make(chan task.Task, p.Total)
		for _, v := range box {
			source <- v
		}
		close(source)

		// generate cancelable context
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// set cancel signals
		sigCh := make(chan os.Signal)
		defer close(sigCh)
		// watch SIGINT, SIGSTOP
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGSTOP)
		go func() {
			<-sigCh
			cancel()
		}()

		// first stage -> write files for simulation
		//first := p.AddStage(x, source, "write", write.Write{
		//	Tmp: config.Templates,
		//})

		// second stage -> simulation with hspice
		second := p.AddStage(x, source, "simulation", simulation.HSPICE{
			Path:    config.HSPICE.Path,
			Options: config.HSPICE.Options,
			Tmp:     config.Templates,
			Log:     log.WithField("at", "simulation").Logger,
		})

		// third stage -> extract with waveview
		third := p.AddStage(y, second, "extract", extract.WaveView{
			Path: config.WaveView.Path,
			Log:  log.WithField("at", "extract").Logger,
		})

		//fourth stage -> push with taa
		fourth := p.AddStage(z, third, "push", push.Taa{
			ConfigFile: config.Taa.ConfigFile,
			TaaPath:    config.Taa.Path,
			Log:        log.WithField("at", "push").Logger,
		})

		// fifth stage -> remove csv, spi
		fifth := p.AddStage(1, fourth, "remove", remove.Remove{})

		// error channel
		errCh := make(chan error)
		defer close(errCh)

		// Start PipeLine
		go func() {
			errCh <- p.Start(ctx)
		}()

		select {
		case err := <-errCh:
			if err != nil {
				logrus.WithError(err).Fatal("Pipeline task was failed")
			}
		}

		spin := spinner.New(spinner.CharSets[14], time.Millisecond*500)
		spin.FinalMSG = "done: clean up"
		spin.Suffix = "clean up"
		spin.Start()

		successes := 0

		var dirs []string
		for v := range fifth {
			var p string
			var err error
			if !keepcsv {
				p, err = filepath.Abs(filepath.Join(v.Files.Directories.DstDir, "..", ".."))
			} else {
				p, err = filepath.Abs(filepath.Join(v.Files.Directories.DstDir))
			}
			if err != nil {
				log.WithError(err).Warn("can not resolve path: ", p)
			}
			dirs = append(dirs, p)
			successes++
		}

		for _, v := range dirs {
			err := remove.Do(context.Background(), v)
			if err != nil {
				log.WithError(err).Error("failed remove dir: ", v)
			}
		}

		spin.Stop()

		fmt.Println()

		msg := "finished: avv run"
		log.Info(msg)
		logrus.Info(msg)

		if slack {
			config.SlackConfig.PostMessage(fmt.Sprintf(":seikou: %d\n:sippai: %d", successes, p.Total-successes))
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().IntP("count", "n", 1, "number of task")
	runCmd.Flags().Bool("all", false, "")

<<<<<<< HEAD
	runCmd.Flags().IntP("simulateParallel", "x", 1, "HSPICEの並列数です")
	runCmd.Flags().IntP("extractParallel", "y", 1, "WaveViewの並列数です")
	runCmd.Flags().IntP("pushParallel", "z", 1, "taaの並列数です")

	runCmd.Flags().IntP("maxRetry", "m", 3, "各ステージの処理が失敗したときに再実行する回数です")
	viper.BindPFlag("MaxRetry", runCmd.Flags().Lookup("maxRetry"))
=======
func (rt RunTask) Run(ctx context.Context) {

	l := logrus.WithField("at", "avv run")
	l.Info(Version)
	l.Info("Start run command")
	l.Info("Number of Task=", len(rt))

	var t []ITask
	for _, r := range rt {
		t = append(t, r)
	}

	begin := time.Now()
	p := NewPipeLine(t)
	res := p.Start(ctx)
	end := time.Now()
	var s, f []Task

	for _, v := range res {
		if v.Status {
			s = append(s, v.Task)
		} else {
			f = append(f, v.Task)
		}
	}

	var sp, fp string
	{
		sp = PathJoin(DoneDir(), time.Now().Format("2006-01-02-15-04-05.json"))
		l.Info("Write Success Tasks to file: ", sp)
		b, err := json.MarshalIndent(&s, "", "  ")
		if err != nil {
			l.WithError(err).Fatal("Failed Marshal Success Tasks")
		}
		err = ioutil.WriteFile(sp, b, 0644)
		if err != nil {
			l.WithError(err).Fatal("Failed Write to ", sp)
		}
	}
	{
		fp = PathJoin(FailedDir(), time.Now().Format("2006-01-02-15-04-05.json"))
		l.Info("Write Failed Tasks to file: ", fp)
		b, err := json.MarshalIndent(&f, "", "  ")
		if err != nil {
			l.WithError(err).Fatal("Failed Marshal Failed Tasks")
		}
		err = ioutil.WriteFile(sp, b, 0644)
		if err != nil {
			l.WithError(err).Fatal("Failed Write to ", fp)
		}
	}

	config.SlackConfig.PostMessage(fmt.Sprintf(":seikou: %d\n:sippai: %d\n開始時間：%s\n終了時間：%s",
		len(s), len(f), begin.Format(time.ANSIC), end.Format(time.ANSIC)))
>>>>>>> master

	runCmd.Flags().Bool("slack", true, "すべてのタスクが終わったときにSlackへ投稿します")
	runCmd.Flags().Bool("keepcsv", false, "CSVを残します")
}
