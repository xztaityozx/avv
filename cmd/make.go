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
	"encoding/json"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xztaityozx/avv/task"
	"io/ioutil"
	"path/filepath"
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
		log.Info("avv make")

		tasks, err := task.Generate(config)
		if err != nil {
			log.Fatal(err)
		}

		for _, v := range tasks {
			b, err := json.Marshal(&v)
			if err != nil {
				log.WithError(err).Error("Failed make task json")
			}

			path := filepath.Join(config.Default.TaskDir(), v.Hash()+".json")
			if err := ioutil.WriteFile(path, b, 0644); err != nil {
				log.WithError(err).Error("Failed write task json")
			}
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

	makeCmd.Flags().IntP("times", "t", 0, "モンテカルロシミュレーション1回当たりの回数")
	viper.BindPFlag("Default.Sweeps", makeCmd.Flags().Lookup("times"))
	viper.BindPFlag("Default.SEED.Start", makeCmd.Flags().Lookup("start"))
	viper.BindPFlag("Default.SEED.End", makeCmd.Flags().Lookup("end"))

	makeCmd.Flags().String("basedir", "", "シミュレーションの結果を書き出す親ディレクトリ")

	viper.BindPFlag("Default.BaseDir", makeCmd.Flags().Lookup("basedir"))

}
