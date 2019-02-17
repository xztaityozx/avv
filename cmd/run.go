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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"time"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "",
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
		fmt.Println(">>> ")

		s := bufio.NewScanner(os.Stdin)
		s.Scan()
		ans := s.Text()
		config.AutoDBBackUp = (ans == "y") || config.AutoDBBackUp

		err := rt.BackUp()
		if err != nil {
			log.WithError(err).Fatal("Failed DB Backup")
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().IntP("count", "n", 1, "number of task")
	runCmd.Flags().Bool("all",false,"")
	runCmd.Flags().String("taskdir", "", "タスクファイルが保存されてるディレクトリへのパスです.省略すると設定ファイルにかかれている値と同じになります")
	viper.BindPFlag("TaskDir", runCmd.Flags().Lookup("taskdir"))
}

func (rt RunTask) Run(ctx context.Context) {
	var tasks []ITask
	for _, v := range rt {
		tasks = append(tasks, SimulationTask{Task: v})
	}

	begin := time.Now()

	s, f, err := PipeLine{}.Start(ctx, tasks,
		// HSPICE Pipe
		Pipe{
			Name:       string(HSPICE),
			Parallel:   config.ParallelConfig.HSPICE,
			RetryLimit: config.RetryConfig.HSPICE,
			Converter: func(task Task) ITask {
				task.Stage = WaveView
				return ExtractTask{Task: task}
			},
		},
		// WaveView Pipe
		Pipe{
			Name:       "Extract",
			Parallel:   config.ParallelConfig.WaveView,
			RetryLimit: config.RetryConfig.WaveView,
			Converter: func(task Task) ITask {
				task.Stage = CountUp
				return CountTask{Task: task}
			},
		},
		// CountUp Pipe
		Pipe{
			Name:       "CountUp",
			Parallel:   config.ParallelConfig.CountUp,
			RetryLimit: config.RetryConfig.CountUp,
			Converter: func(task Task) ITask {
				task.Stage = DBAccess
				return DBAccessTask{Task: task}
			},
		},
		// DB Access Pipe
		Pipe{
			Name:       "DB Access",
			Parallel:   config.ParallelConfig.DB,
			RetryLimit: config.ParallelConfig.DB,
			Converter: func(task Task) ITask {
				task.Stage = "Finish"
				return nil
			},
		})

	length := len(s)
	if len(f) > length {
		length = len(f)
	}

	fmt.Println("Success\tFailed")

	for i := 0; i < length; i++ {
		if i < len(s) {
			fmt.Print(s[i], "\t")
		}
		if i < len(f) {
			fmt.Print(f[i])
		}

		fmt.Println()
	}

	if err != nil {
		log.WithError(err).Error("Failed run command")
	}

	end := time.Now()

	config.SlackConfig.PostMessage(fmt.Sprintf(":seikou: %d\n:sippai: %d\n開始時間: %s\n終了時間: %s",
		len(s), len(f),
		begin.Format("2006/01/02/15:04:05"),
		end.Format("2006/01/02/15:04:05")))
}
