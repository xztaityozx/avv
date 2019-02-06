package cmd

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"os/exec"
	"strings"
)

var extractCmd = &cobra.Command{
	Use:     "extract",
	Aliases: []string{"ext"},
	Short:   `WaveViewの出力したcsvから必要なところだけ取り出します`,
	Long:    ``,
	PreRun: func(cmd *cobra.Command, args []string) {
		if len(args) < 3 {
			log.WithField("command", "extract").Fatal("引数が足りません")
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		f := Filter{
			SignalName: args[0],
		}

		for _, v := range strings.Split(args[1], ",") {
			f.Status = append(f.Status, v)
		}

		dst, _ := os.Getwd()
		if len(args) == 4 {
			dst = args[3]
		}

		p, err := f.MakeResultFileFromCSV(args[2], dst, config.Default.PlotPoint.Count())
		if err != nil {
			log.WithField("command", "extract").Fatal(err)
		}

		log.WithField("command", "extract").Info("Write To: ", p)
	},
}

func init() {
	rootCmd.AddCommand(extractCmd)

	extractCmd.Flags().Float64("start", 0, "プロットの開始時間")
	extractCmd.Flags().Float64("step", 0, "プロットの刻み幅")
	extractCmd.Flags().Float64("stop", 0, "プロットの終了時間")

	viper.BindPFlag("Default.PlotPoint.Start", extractCmd.Flags().Lookup("start"))
	viper.BindPFlag("Default.PlotPoint.Step", extractCmd.Flags().Lookup("step"))
	viper.BindPFlag("Default.PlotPoint.Stop", extractCmd.Flags().Lookup("stop"))
}

type ExtractTask struct {
	Task Task
}

// Generate wv Command
func (wvc WaveViewConfig) GetCommand(ace string) string {
	return fmt.Sprintf("%s -k -ace_no_gui %s &> wv.log", wvc.Command, ace)
}

// Generate extract command
func (t ExtractTask) GetExtractCommand() string {
	var rt []string

	// append cd command
	rt = append(rt, t.Task.GetCdCommand())
	// append waveview command
	rt = append(rt, config.WaveView.GetCommand(t.Task.SimulationFiles.ACEScript))

	return strings.Join(rt, " ")
}

func (t ExtractTask) String() string {
	return fmt.Sprint()
}

// Run extract command with WaveView
// returns: errors
func (t ExtractTask) Run(parent context.Context) Result {

	l := log.WithField("at", "ExtractTask")

	// Make AddFile
	t.Task.SimulationFiles.AddFile.Make(t.Task.SimulationDirectories.BaseDir)

	// Make results.xml
	if path, err := t.Task.MakeResultsXml(); err != nil {
		l.WithError(err).Error("Failed make results.xml")
		return Result{
			Task:   t.Task,
			Status: false,
		}
	} else {
		t.Task.SimulationFiles.ResultsXML = path
	}

	// Make resultsMap.xml
	if path, err := t.Task.MakeMapXml(); err != nil {
		l.WithError(err).Error("Failed make resultsMap.xml")
		return Result{
			Task:   t.Task,
			Status: false,
		}
	} else {
		t.Task.SimulationFiles.ResultsMapXML = path
	}

	cmdStr := t.GetExtractCommand()
	command := exec.Command("bash", "-c", cmdStr)
	_ , err := command.CombinedOutput()
	if err != nil {
		logfile := PathJoin(t.Task.SimulationDirectories.DstDir, "wv.log")
		l.WithError(err).Error("Failed Extract: " + FU.Cat(logfile))
		return Result{
			Task:   t.Task,
			Status: false,
		}
	}

	for _, v := range t.Task.PlotPoint.Filters {
		// csvをそれぞれの信号線のディレクトリに移動させる
		oldPath := PathJoin(t.Task.SimulationDirectories.DstDir, v.SignalName+".csv")
		FU.TryMkDir(PathJoin(t.Task.SimulationDirectories.ResultDir, v.SignalName))
		newPath := PathJoin(t.Task.SimulationDirectories.ResultDir, v.SignalName, fmt.Sprintf("SEED%05d.csv", t.Task.SEED))

		if err := os.Rename(oldPath, newPath); err != nil {
			l.WithError(err).Error("Failed Rename ", oldPath, " to ", newPath)
			return Result{
				Task:   t.Task,
				Status: false,
			}
		}
	}

	return Result{
		Task:   t.Task,
		Status: true,
	}
}

func (t ExtractTask) Self() Task {
	return t.Task
}

// Make result file from csv file which generated by WaveView
// returns: file path, error
func (f Filter) MakeResultFileFromCSV(srcDir, dstDir string, plotStepCount int) (string, error) {
	dstFile := PathJoin(dstDir, f.SignalName+".csv")
	srcFile := PathJoin(srcDir, f.SignalName+".csv")

	fp, err := os.Open(srcFile)
	defer fp.Close()
	if err != nil {
		return "", errors.New("MakeResultFileFromCSV: can not found file " + srcFile)
	}

	scan := bufio.NewScanner(fp)
	var box []string
	var line []string
	for scan.Scan() {
		v := strings.Trim(scan.Text(), "\n")
		if len(v) == 0 || v[0] == '#' || v[0] == 'T' {
			continue
		}
		col := strings.Split(strings.Replace(v, " ", "", -1), ",")
		if len(col) < 2 {
			return "", errors.New("Invalid File Format\n")
		}

		line = append(line, v)
		if len(line) == plotStepCount {
			box = append(box, strings.Join(line, ","))
			line = []string{}
		}
	}

	FU.WriteSlice(dstFile, box, "\n")

	return dstFile, nil
}
