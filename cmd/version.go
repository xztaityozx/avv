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
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print avv version",
	Long:  `バージョン情報を出力して終了します`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(Version)
	},
}

type AVVVersion struct {
	Major  int
	Minor  int
	Build  int
	Date   string
	Status string
}

var Version = AVVVersion{
	Major:  0,
	Minor:  1,
	Build:  25,
	Date:   "2019/02/20",
	Status: "Development",
}

func (av AVVVersion) String() string {
	return fmt.Sprintf("avv v%d.%d.%d %s %s\n\nAuthor: xztaityozx\nRepository: https://github.com/xztaityozx/avv\nLicense: MIT",
		av.Major, av.Minor, av.Build, av.Date, av.Status)
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
