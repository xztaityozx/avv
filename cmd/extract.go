package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"os/exec"
	"strconv"
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
			value, err := strconv.ParseFloat(v, 64)
			if err != nil {
				log.WithField("command", "extract").Fatal(err)
			}
			f.Values = append(f.Values, value)
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

// Generate wv Command
func (wvc WaveViewConfig) GetCommand(ace string) string {
	return fmt.Sprintf("%s -k -ace_no_gui %s &> wv.log", wvc.Command, ace)
}

// Generate extract command
func (t Task) GetExtractCommand() string {
	var rt []string

	// append cd command
	rt = append(rt, t.GetCdCommand())
	// append waveview command
	rt = append(rt, config.WaveView.GetCommand(t.SimulationFiles.ACEScript))

	return strings.Join(rt, " ")
}

// Run extract command with WaveView
// returns: csv file, errors
func (t Task) RunExtract() ([]string, error) {
	// Make AddFile
	t.SimulationFiles.AddFile.Make(t.SimulationDirectories.BaseDir)

	// Make results.xml
	if path, err := t.MakeResultsXml(); err != nil {
		return nil, err
	} else {
		t.SimulationFiles.ResultsXML = path
	}

	// Make resultsMap.xml
	if path, err := t.MakeMapXml(); err != nil {
		return nil, err
	} else {
		t.SimulationFiles.ResultsMapXML = path
	}

	cmdStr := t.GetExtractCommand()
	command := exec.Command("bash", "-c", cmdStr)
	out, err := command.CombinedOutput()
	if err != nil {
		logfile := PathJoin(t.SimulationDirectories.DstDir, "wv.log")
		return nil, errors.New("Failed Extract: " + FU.Cat(logfile))
	} else {
		log.WithField("at", "Task.Run").Info(string(out))
	}

	var rtPath []string
	for _, v := range t.PlotPoint.Filters {
		rtPath = append(rtPath,
			PathJoin(t.SimulationDirectories.DstDir, fmt.Sprintf("%s.csv", v.SignalName)))
	}

	return rtPath, nil
}

// Make result file from csv file which generated by WaveView
// returns: file path, error
func (f Filter) MakeResultFileFromCSV(srcDir, dstDir string, plotStepCount int) (string, error) {
	dstFile := PathJoin(dstDir, f.SignalName+".csv")
	srcFile := PathJoin(srcDir, f.SignalName+".csv")

	fp, err := os.Open(srcFile)
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
