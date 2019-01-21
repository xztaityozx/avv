package cmd

import (
	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestAllTask(t *testing.T) {

	home, _ := homedir.Dir()
	config.DstDir = PathJoin(home, "Dst")
	config.SrcDir = PathJoin(home, "Src")

	as := assert.New(t)
	t.Run("001 NewTask", func(t *testing.T) {
		expect := Task{
			ParallelConfig: ParallelConfig{
				CountUp:  10,
				WaveView: 20,
				HSPICE:   30,
			},
			PlotPoint: PlotPoint{
				SignalNames: []string{"A", "B", "C"},
				Stop:        17.5,
				Step:        7.5,
				Start:       2.5,
			},
			SEED: SEED{
				Start: 1,
				End:   2000,
			},
			TaskNameFileName: "TaskName",
			SrcDir:           "Src",
			DstDir:           "Dst",
			Vtp: Transistor{
				Sigma:     0.7,
				Deviation: 1.0,
				Threshold: 0.2,
			},
			Vtn: Transistor{
				Threshold: 0.3,
				Deviation: 1.2,
				Sigma:     0.3,
			},
			AutoRemove: true,
		}

		config.ParallelConfig = ParallelConfig{
			CountUp:  10,
			WaveView: 20,
			HSPICE:   30,
		}

		actual := NewTask("Dst", "Src", "TaskName", 1, 2000, 0.3, 0.3,
			1.2, 0.2, 0.7, 1.0, true)

		as.Equal(expect, actual)

	})

	t.Run("002 Task_MkDirs", func(t *testing.T) {
		task := NewTask("Dst", "Src", "TaskName", 1, 20, 0.3, 0.3,
			1.2, 0.2, 0.7, 1.0, true)

		task.MkDirs()

		os.RemoveAll(config.DstDir)
	})
}
