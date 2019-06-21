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
	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/spf13/viper"
	"github.com/xztaityozx/avv/remove"
	"time"

	"github.com/spf13/cobra"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"rm"},
	Short:   "",
	Long:    `remove temporally simulation files`,
	Run: func(cmd *cobra.Command, args []string) {
		base := config.Default.BaseDir

		s := spinner.New(spinner.CharSets[14], time.Millisecond*500)
		s.FinalMSG = color.New(color.FgGreen).Sprint("Finished")
		s.Suffix = color.New(color.FgHiYellow).Sprint("processing...")
		s.Start()
		defer s.Stop()

		err := remove.Do(context.Background(), base)
		if err != nil {
			log.Fatal(err)
		}

	},
}

func init() {
	rootCmd.AddCommand(removeCmd)

	removeCmd.Flags().String("target", "", "/path/to/basedir")
	viper.BindPFlag("Default.BaseDir", removeCmd.Flags().Lookup("target"))
}