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
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "シミュレーションを実行します",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var rt RunTask
		if len(args) == 0 {
			cnt, _ := cmd.Flags().GetInt("count")
			all, _ := cmd.Flags().GetBool("all")
			if all {
				cnt = -1
			}
			rt = GetTasksFromTaskDir(cnt)
		} else {
			rt = GetTasksFromFiles(args...)
		}

		fmt.Println("全てのDBのバックアップを取りますか？")
		fmt.Println("[y]:はい [n]:いいえ")
		fmt.Print(">>> ")

		s := bufio.NewScanner(os.Stdin)
		s.Scan()
		ans := s.Text()
		config.AutoDBBackUp = (ans == "y") || config.AutoDBBackUp

		err := rt.BackUp()
		if err != nil {
			log.WithError(err).Fatal("Failed DB Backup")
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		sigCh := make(chan os.Signal)
		defer close(sigCh)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGSTOP, syscall.SIGSTOP)
		go func() {
			<-sigCh
			cancel()
		}()

		ch := make(chan struct{})
		defer close(ch)

		go func() {
			rt.Run(ctx)
			ch <- struct{}{}
		}()

		select {
		case <-ctx.Done():
		case <-ch:
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().IntP("count", "n", 1, "number of task")
	runCmd.Flags().Bool("all", false, "")
	runCmd.Flags().String("taskdir", "", "タスクファイルが保存されてるディレクトリへのパスです.省略すると設定ファイルにかかれている値と同じになります")
	viper.BindPFlag("TaskDir", runCmd.Flags().Lookup("taskdir"))
}

func (rt RunTask) Run(ctx context.Context) {

	l := logrus.WithField("at", "avv run")
	l.Info(Version)
	l.Info("Start run command")
	l.Info("Number of Task=", len(rt))

	var tasks []ITask
	{
		m := make(map[int64]bool)
		for _, v := range rt {
			tasks = append(tasks, SimulationTask{Task: v})
			m[v.TaskId] = true
		}

		logrus.Info("いくつかのTaskIdを発見しました。結果へアクセスする際に使うので控えておいてね")
		for i := range m {
			fmt.Println(i)
		}
	}

	begin := time.Now()

	s, f, err := PipeLine{}.Start(ctx, tasks,
		// HSPICE Pipe
		Pipe{
			Name:       string(HSPICE),
			Parallel:   config.ParallelConfig.HSPICE,
			RetryLimit: config.RetryConfig.HSPICE,
			AutoRetry:  true,
			Converter: func(task Task) ITask {
				task.Stage = WaveView
				return ExtractTask{Task: task}
			},
			FailedConverter: func(task Task) ITask {
				return SimulationTask{Task:task}
			},
		},
		// WaveView Pipe
		Pipe{
			Name:       "Extract",
			Parallel:   config.ParallelConfig.WaveView,
			RetryLimit: config.RetryConfig.WaveView,
			AutoRetry:  true,
			Converter: func(task Task) ITask {
				task.Stage = CountUp
				return CountTask{Task: task}
			},
			FailedConverter: func(task Task) ITask {
				return ExtractTask{Task:task}
			},
		},
		// CountUp Pipe
		Pipe{
			Name:       "CountUp",
			Parallel:   config.ParallelConfig.CountUp,
			RetryLimit: config.RetryConfig.CountUp,
			AutoRetry:  true,
			Converter: func(task Task) ITask {
				task.Stage = DBAccess
				return DBAccessTask{Task: task}
			},
			FailedConverter: func(task Task) ITask {
				return CountTask{Task:task}
			},
		},
		// DB Access Pipe
		Pipe{
			Name:       "DB Access",
			Parallel:   config.ParallelConfig.DB,
			RetryLimit: config.RetryConfig.DB,
			AutoRetry:  true,
			Converter: func(task Task) ITask {
				task.Stage = "Remove"
				return RemoveTask{
					Task: task,
				}
			},
			FailedConverter: func(task Task) ITask {
				return DBAccessTask{Task:task}
			},
		},
		// Remove Pipe
		Pipe{
			Name:       "Remove",
			Parallel:   config.ParallelConfig.Remove,
			RetryLimit: 0,
			AutoRetry:  false,
			Converter: func(task Task) ITask {
				task.Stage = "Finish"
				return task
			},
			FailedConverter: func(task Task) ITask {
				return RemoveTask{Task:task}
			},
		})

	if err != nil {
		l.WithError(err).Error("Failed run command")
	}

	end := time.Now()

	config.SlackConfig.PostMessage(fmt.Sprintf(":seikou: %d\n:sippai: %d\n開始時間: %s\n終了時間: %s",
		len(s), len(f),
		begin.Format("2006/01/02/15:04:05"),
		end.Format("2006/01/02/15:04:05")))
}
