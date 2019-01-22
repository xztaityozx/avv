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
	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

// makeCmd represents the make command
var makeCmd = &cobra.Command{
	Use:   "make",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info(config)
	},
}

func init() {
	rootCmd.AddCommand(makeCmd)

	makeCmd.Flags().Float64("PlotStart", 2.5, "プロットの始点[ns]")
	makeCmd.Flags().Float64("PlotStep", 7.5, "プロットの刻み幅[ns]")
	makeCmd.Flags().Float64("PlotStop", 17.5, "プロットの終点[ns]")
	makeCmd.Flags().StringSlice("signals", []string{"N1", "N2", "BLB", "BL"}, "プロットしたい信号名")

	viper.BindPFlag("Default.PlotPoint.Start", makeCmd.Flags().Lookup("PlotStart"))
	viper.BindPFlag("Default.PlotPoint.Step", makeCmd.Flags().Lookup("PlotStep"))
	viper.BindPFlag("Default.PlotPoint.Stop", makeCmd.Flags().Lookup("PlotStop"))
	viper.BindPFlag("Default.PlotPoint.SignalNames", makeCmd.Flags().Lookup("signals"))

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

}

type MakeRequest struct {
	Task Task
}
