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
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

// makeCmd represents the make command
var makeCmd = &cobra.Command{
	Use:   "make",
	Short: "タスクを作ります",
	Long: `パラメータを指定してタスクファイルを生成します
SEEDごとに1つのファイルが生成されます
生成先は設定ファイルの "TaskDir" です

指定しなかった値は設定ファイルの値が使われます


`,
	Run: func(cmd *cobra.Command, args []string) {

		mr := MakeRequest{
			Task: NewTask(),
			SEED: SEED{
				Start: config.DefaultSEEDRange.Start,
				End:   config.DefaultSEEDRange.End},
			TaskDir: config.TaskDir,
		}

		// cancel付きcontext
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// SIGNALをトラップする
		sigCh := make(chan os.Signal)
		defer close(sigCh)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGSTOP, syscall.SIGQUIT)

		go func() {
			<-sigCh
			cancel()
		}()

		// 終了通知チャンネル
		ch := make(chan struct{})
		defer close(ch)

		// DBにアクセスするのでちょっと重い
		go func() {
			clo := func() { ch <- struct{}{} }
			defer clo()
			if err := mr.MakeTaskFiles(ctx); err != nil {
				log.Fatal("MakeTaskFiles: ", err)
			}
			logrus.Info("Task File Wrote to ", ReserveDir())
		}()

		// 待機
		select {
		case <-ctx.Done():
		case <-ch:
		}

	},
}

func init() {
	rootCmd.AddCommand(makeCmd)

	makeCmd.Flags().Float64P("PlotStart", "a", 2.5, "プロットの始点[ns]")
	makeCmd.Flags().Float64P("PlotStep", "b", 7.5, "プロットの刻み幅[ns]")
	makeCmd.Flags().Float64P("PlotStop", "c", 17.5, "プロットの終点[ns]")

	viper.BindPFlag("Default.PlotPoint.Start", makeCmd.Flags().Lookup("PlotStart"))
	viper.BindPFlag("Default.PlotPoint.Step", makeCmd.Flags().Lookup("PlotStep"))
	viper.BindPFlag("Default.PlotPoint.Stop", makeCmd.Flags().Lookup("PlotStop"))

	makeCmd.Flags().Float64P("VtpVoltage", "P", 0, "Vtpのしきい値電圧")
	makeCmd.Flags().Float64P("VtnVoltage", "N", 0, "Vtnのしきい値電圧")
	makeCmd.Flags().Float64("VtpDeviation", 1.0, "Vtpの偏差")
	makeCmd.Flags().Float64("VtnDeviation", 1.0, "Vtnの偏差")
	makeCmd.Flags().Float64P("sigma", "S", 0.046, "VtpとVtnのシグマ")
	makeCmd.Flags().Float64("VtpSigma", 0, "Vtpのシグマ")
	makeCmd.Flags().Float64("VtnSigma", 0, "Vtnのシグマ")

	viper.BindPFlag("Default.Vtn.Threshold", makeCmd.Flags().Lookup("VtnVoltage"))
	viper.BindPFlag("Default.Vtn.Deviation", makeCmd.Flags().Lookup("VtnDeviation"))
	viper.BindPFlag("Default.Vtn.Sigma", makeCmd.Flags().Lookup("VtnSigma"))
	viper.BindPFlag("Default.Vtp.Threshold", makeCmd.Flags().Lookup("VtpVoltage"))
	viper.BindPFlag("Default.Vtp.Deviation", makeCmd.Flags().Lookup("VtpDeviation"))
	viper.BindPFlag("Default.Vtp.Sigma", makeCmd.Flags().Lookup("VtpSigma"))

	makeCmd.Flags().BoolP("autoremove", "r", false, "波形データを自動で削除します")
	makeCmd.Flags().Int("start", 0, "SEEDの始点")
	makeCmd.Flags().Int("end", 0, "SEEDの終点")

	viper.BindPFlag("Default.AutoRemove", makeCmd.Flags().Lookup("autoremove"))
	viper.BindPFlag("DefaultSEEDRange.Start", makeCmd.Flags().Lookup("start"))
	viper.BindPFlag("DefaultSEEDRange.End", makeCmd.Flags().Lookup("end"))

	makeCmd.Flags().IntP("times", "t", 0, "モンテカルロシミュレーション1回当たりの回数")
	viper.BindPFlag("Default.Times", makeCmd.Flags().Lookup("times"))

	makeCmd.Flags().String("basedir", "", "シミュレーションの結果を書き出す親ディレクトリ")
	makeCmd.Flags().String("logdir", "", "ログを格納するディレクトリ")

	viper.BindPFlag("Default.SimulationDirectories.BaseDir", makeCmd.Flags().Lookup("basedir"))
	viper.BindPFlag("LogDir", makeCmd.Flags().Lookup("logdir"))

	makeCmd.Flags().StringP("DB", "d", "", "path to Output DataBase")
	viper.BindPFlag("Default.Repository.Path", makeCmd.Flags().Lookup("DB"))

}

// MakeRequest
type MakeRequest struct {
	Task    Task
	SEED    SEED
	TaskDir string
}

// make Task file
// output task.json to [TaskDir/reserve]
// return: error
func (m MakeRequest) MakeTaskFiles(ctx context.Context) error {

	for seed := m.SEED.Start; seed <= m.SEED.End; seed++ {
		data := m.Task
		data.SEED = seed
		data.Stage = HSPICE
		data.SimulationFiles.AddFile.SEED = seed

		if b, err := json.Marshal(data); err != nil {
			return err
		} else {
			// 書き出し先
			path := ReserveDir()
			FU.TryMkDir(path)

			// ファイル名[時間]-[Vtn]-[VtnSigma]-[Vtp]-[VtpSigma]-[回数]-[SEED].json
			path = PathJoin(path, fmt.Sprintf("%s-Vtn%.4f-Sigma%.4f-Vtp%.4f-Sigma%.4f-Times%05d-SEED%05d.json",
				time.Now().Format("20060102150405"),
				data.Vtn.Threshold, data.Vtn.Sigma,
				data.Vtp.Threshold, data.Vtp.Sigma,
				data.Times,
				data.SEED))
			err := ioutil.WriteFile(path, b, 0644)
			if err != nil {
				return err
			}
		}
	}
	log.Info("Make Request: Write ", m.SEED.End, "task files")
	return nil
}
