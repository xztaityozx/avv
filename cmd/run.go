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
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	var t []ITask
	for _, r := range rt {
		t = append(t, r)
	}

	begin := time.Now()
	p := NewPipeLine(t)
	res := p.Start(ctx)
	end := time.Now()
	var s,f []Task

	for _, v := range res {
		if v.Status {
			s=append(s, v.Task)
		} else {
			f = append(f, v.Task)
		}
	}


	var sp,fp string
	{
		sp = PathJoin(DoneDir(),time.Now().Format("2006-01-02-15-04-05.json"))
		l.Info("Write Success Tasks to file: ",sp)
		b, err := json.MarshalIndent(&s, "","  ")
		if err != nil {
			l.WithError(err).Fatal("Failed Marshal Success Tasks")
		}
		err = ioutil.WriteFile(sp, b, 0644)
		if err != nil {
			l.WithError(err).Fatal("Failed Write to ", sp)
		}
	}
	{
		fp = PathJoin(FailedDir(),time.Now().Format("2006-01-02-15-04-05.json"))
		l.Info("Write Failed Tasks to file: ", fp)
		b, err := json.MarshalIndent(&f, "", "  ")
		if err != nil {
			l.WithError(err).Fatal("Failed Marshal Failed Tasks")
		}
		err = ioutil.WriteFile(sp, b, 0644)
		if err != nil {
			l.WithError(err).Fatal("Failed Write to ",fp)
		}
	}

	config.SlackConfig.PostMessage(fmt.Sprintf(":seikou: %d\n:sippai: %d\n開始時間：%s\n終了時間：%s",
		len(s),len(f),begin.Format(time.ANSIC),end.Format(time.ANSIC)))

}
