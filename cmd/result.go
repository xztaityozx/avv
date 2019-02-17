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
	"github.com/briandowns/spinner"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

// resultCmd represents the result command
var resultCmd = &cobra.Command{
	Use:   "result",
	Aliases:[]string{"db","rst"},
	Short: "",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		var rr ResultRecord
		rr.TaskId, _ = cmd.Flags().GetInt64("TaskId")
		DBSearch(rr)
	},
}

var dbTransistorCmd = &cobra.Command{
	Use:"tran",
	Aliases:[]string{"tr"},
	Short:"トランジスタIDを検索します",
	Run: func(cmd *cobra.Command, args []string) {
		var t Transistor
		t.Deviation, _ = cmd.Flags().GetFloat64("Deviation")
		t.Sigma,_=cmd.Flags().GetFloat64("Sigma")
		t.Threshold,_=cmd.Flags().GetFloat64("Threshold")

		DBSearch(t)
	},
}

var dbParameterCmd = &cobra.Command{
	Use:"param",
	Aliases:[]string{"pa"},
	Short:"パラメータIDを検索します",
	Run: func(cmd *cobra.Command, args []string) {
		p := Parameter{}
		p.Times, _ = cmd.Flags().GetInt64("Times")
		p.VtnId,_=cmd.Flags().GetInt64("VtnId")
		p.VtpId,_=cmd.Flags().GetInt64("VtpId")

		DBSearch(p)
	},
}

var dbTaskGroupCmd = &cobra.Command{
	Use:"task",
	Aliases:[]string{"t"},
	Short:"タスクIDを検索します",
	Run: func(cmd *cobra.Command, args []string) {
		t := TaskGroup{}
		t.ParamsId,_=cmd.Flags().GetInt64("ParamsId")
		t.SeedEnd,_=cmd.Flags().GetInt64("SeedEnd")
		t.SeedStart,_=cmd.Flags().GetInt64("SeedStart")

		DBSearch(t)
	},
}

var DBtoJson bool
var DBtoOutFile string

func init() {
	rootCmd.AddCommand(resultCmd)
	resultCmd.AddCommand(dbTransistorCmd)
	resultCmd.AddCommand(dbParameterCmd)
	resultCmd.AddCommand(dbTaskGroupCmd)

	// Transistor Query
	dbTransistorCmd.Flags().Float64("Threshold", -1, "検索したいトランジスタのしきい値です")
	dbTransistorCmd.Flags().Float64("Deviation", -1, "検索したいトランジスタの偏差です")
	dbTransistorCmd.Flags().Float64("Sigma",-1,"検索したいトランジスタのシグマです")

	// Parameter
	dbParameterCmd.Flags().Int64("VtnID", 0, "検索したいパラメータIDのVtnのIDです")
	dbParameterCmd.Flags().Int64("VtpID", 0, "検索したいパラメータIDのVtpのIDです")
	dbParameterCmd.Flags().Int64("Times",0,"検索したいパラメータIDのMCSの回数です")

	// TaskGroup
	dbTaskGroupCmd.Flags().Int64("ParamsId",0,"検索したいタスクIDのパラメータのIDです")
	dbTaskGroupCmd.Flags().Int64("SeedStart",0,"検索したいタスクIDのシードの初期値です")
	dbTaskGroupCmd.Flags().Int64("SeedEnd",0,"検索したいタスクIDのシードの終了値です")

	// Result Record
	resultCmd.Flags().Int64("TaskId", 0, "検索したい結果のタスクIDです")

	resultCmd.PersistentFlags().String("db","","アクセスするDBへのパスです。指定しないならコンフィグのものが使われます")
	viper.BindPFlag("Default.Repository.Path",resultCmd.PersistentFlags().Lookup("db"))
	resultCmd.PersistentFlags().BoolVar(&DBtoJson, "json", false,"結果をJSON形式で出力します")
	resultCmd.PersistentFlags().StringVarP(&DBtoOutFile,"out","o","","出力先のファイルです。指定しない場合STDOUTに出力されます")
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

func  DBSearch(t IDBQuery){
	db,err := config.Default.Repository.Connect()
	defer db.Db.Close()
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal)
	defer close(sigCh)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGSTOP, syscall.SIGKILL)
	go func() {
		<-sigCh
		cancel()
	}()
	fns := make(chan struct{})
	defer close(fns)

	spin := spinner.New(spinner.CharSets[14],time.Millisecond*500)
	spin.Suffix = "Accessing to "+config.Default.Repository.Path
	spin.FinalMSG = ""
	spin.Writer=os.Stderr
	spin.Start()

	var box []interface{}
	go func() {
		q,p := t.DBQuery()
		_, err = db.WithContext(ctx).Select(&box,q,p...)
		if err != nil {
			log.Fatal(err)
		}
		fns<- struct{}{}
	}()

	select{
	case <-ctx.Done():
	case <-fns:
	}

	spin.Stop()
	l:=logrus.New()
	l.SetOutput(os.Stderr)
	l.Info("Finished Access to DB")

	PrintTo(box)
}

func PrintTo(box []interface{}) {
	text := ""
	var b []byte
	var err error
	if DBtoJson {
		b, err = json.MarshalIndent(box,"","  ")
		if err != nil {
			log.Fatal(err)
		}
		text = string(b)
	} else {
		var x []string
		for _, v := range box {
			x=append(x,fmt.Sprint(v))
		}
		text = strings.Join(x,"\n")
	}
	if len(DBtoOutFile) == 0 {
		fmt.Println(text)
	} else {
		err = ioutil.WriteFile(DBtoOutFile, b, 0644)
		log.Fatal(err)
	}
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
	if p.Times != 0 {
		when = append(when, "Times = ?")
		param=append(param, p.Times)
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

