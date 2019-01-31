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
	"github.com/spf13/cobra"
)

// countCmd represents the count command
var countCmd = &cobra.Command{
	Use:     "count",
	Aliases: []string{"cnt"},
	Short:   "",
	Long: `avv count [file string] [queries strings]

- [query string]
各項の条件式をカンマ区切りで指定できます．条件はすべて && で連結されます

ex)
	avv count file.csv '>=0.4' '>=0.4' '>=0.4'
	
`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.WithField("command", "count").Fatal("引数が足りません")
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		//out, _ := cmd.Flags().GetString("out")
		//wd, _ := os.Getwd()
		//
		//res :=Filter{
		//	SignalName:filepath.Base(wd),
		//	Status:args,
		//}.CountUp(out)
		//
		//fmt.Printf("Total: %d Failure: %d\n",res.Lines,res.Failure)
	},
}

func init() {
	rootCmd.AddCommand(countCmd)

	countCmd.Flags().StringP("out", "o", "result", "書きだすファイルです")
}

//func (f Filter) CountUp(dst string) CountResult {
//	l := log.WithField("at","Filter.CountUp")
//
//}

type CountResult struct {
	Lines   int
	Failure int
}
