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
	"strings"
)

// resultCmd represents the result command
var resultCmd = &cobra.Command{
	Use:   "result",
	Aliases:[]string{"db","rst"},
	Short: "",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("result called")
	},
}

func init() {
	rootCmd.AddCommand(resultCmd)
}

type (
	IDBQuery interface {
		DBQuery() (string, []interface{})
	}
)

func (t Transistor) DBQuery() (string, []interface{}) {
	var param []interface{}
	base := "select * from Transistor"
	var when []string
	if t.Threshold!= -1 {
		when = append(when, "Threshold = ?")
		param = append(param, t.Threshold)
	}
	if t.Sigma != -1 {
		when = append(when, "Sigma = ?")
		param = append(param, t.Sigma)
	}
	if t.Deviation != -1 {
		when = append(when, "Deviation = ?")
		param = append(param, t.Deviation)
	}

	if len(when) != 0 {
		return base + " when " + strings.Join(when, " and "), param
	}
	return base, param
}

func (p Parameter) DBQuery() (string, []interface{}) {
	var param []interface{}
	base := "select * from Parameter"
	var when []string
	if p.VtpId != 0 {
		when = append(when, "VtpId = ?")
		param = append(param, p.VtpId)
	}
	if p.VtnId != 0 {
		when = append(when, "VtnId = ?")
		param = append(param, p.VtnId)
	}
	if p.ParamsId != 0 {
		when = append(when, "ParamsId = ?")
		param=append(param, p.ParamsId)
	}

	if len(base) != 0 {
		return base + " when "+strings.Join(when," and "), param
	}

	return base, param
}

func (rr ResultRecord) DBQuery() (string, []interface{}) {
	var param []interface{}
	var when []string
	base := "select * from Results"
	if rr.TaskId != 0 {
		when = append(when, "TaskId = ?")
		param = append(param, rr.TaskId)
	}
	if rr.Seed != 0 {
		when = append(when, "Seed = ?")
		param = append(param, rr.Seed)
	}
	if rr.Failure != 0 {
		when = append(when, "Failure = ?")
		param = append(param, rr.Failure)
	}

	if len(when) != 0 {
		base += " when " + strings.Join(when, " and ")
	}

	return base, param
}

func (t TaskGroup) DBQuery() (string, []interface{}) {
	var param []interface{}
	var when []string
	base := "select * from Group"
	if t.ParamsId != 0 {
		when=append(when, "ParamsId = ?")
		param=append(param, t.ParamsId)
	}
	if t.SeedStart != 0 {
		when=append(when, "SeedStart = ?")
		param=append(param, t.SeedStart)
	}
	if t.SeedEnd != 0 {
		when=append(when, "SeedEnd = ?")
		param=append(param, t.SeedEnd)
	}

	if len(when) != 0 {
		base += " when "+strings.Join(when," and ")
	}
	return base, param
}

