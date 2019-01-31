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
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
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
		out, _ := cmd.Flags().GetString("out")
		wd, _ := os.Getwd()

		res := Filter{
			SignalName: filepath.Base(wd),
			Status:     args,
		}.CountUp(out)

		cum, _ := cmd.Flags().GetBool("Cumulative")
		if cum {

			// 累積和
			lines := 0
			failure := 0
			for _, v := range res {
				lines += v.Lines
				failure += v.Failure
				fmt.Println(lines, failure)
			}
		} else {
			// そのまま
			for _, v := range res {
				fmt.Println(v.Lines, v.Failure)
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(countCmd)

	countCmd.Flags().StringP("out", "o", "result", "書きだすディレクトリです")
	countCmd.Flags().BoolP("Cumulative", "c", false, "累積和で出力します")
}

func (f Filter) CountUp(resultDir string) []CountResult {
	l := log.WithField("at", "Filter.CountUp")

	{
		//  awk command check
		command := exec.Command("awk", `BEGIN{print "test"}`)
		err := command.Run()
		if err != nil {
			l.Fatal(err)
		}
	}

	for i, v := range f.Status {
		f.Status[i] = fmt.Sprintf("$%d%s", i+1, v)
	}

	// awkのルール
	rule := strings.Join(f.Status, "&&")

	// 並列数え上げ用Worker
	worker := func(files []string) <-chan CountResult {
		receiver := make(chan CountResult, config.ParallelConfig.CountUp)
		for i, v := range files {
			go func(i int, p string) {
				command := exec.Command("awk", `-f,`, p, fmt.Sprintf("BEGIN{s=0}%s{s++}END{print NR, s}", rule))
				b, err := command.CombinedOutput()
				if err != nil {
					l.Fatal(string(b))
				}
				res := strings.Split(string(b), " ")
				lines, err := strconv.Atoi(res[0])
				if err != nil {
					l.Fatal(err)
				}
				failure, err := strconv.Atoi(res[1])
				if err != nil {
					l.Fatal(err)
				}

				receiver <- CountResult{
					SEED:    i + 1,
					Failure: failure,
					Lines:   lines,
				}
			}(i, v)
		}
		return receiver
	}

	// ファイルのListup
	files := FU.Ls(PathJoin(resultDir, f.SignalName))
	var paths []string
	for _, v := range files {
		paths = append(paths, PathJoin(resultDir, f.SignalName, v.Name()))
	}

	rec := worker(paths)
	result := make([]CountResult, len(paths))

	// 並列化
	for i := 0; i < len(paths); i++ {
		cr := <-rec
		result[cr.SEED] = cr
	}

	return result
}

type CountResult struct {
	SEED    int
	Lines   int
	Failure int
}
