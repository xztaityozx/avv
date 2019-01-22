package cmd

import (
	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestAllTask(t *testing.T) {

	home, _ := homedir.Dir()
	config.Default.SimulationDirectories.BaseDir = PathJoin(home, "Base")
	config.Default.SimulationDirectories.NetListDir = PathJoin(home, "NetList")
	config.Default.SEED = 100
	config.Default.Vtn = Transistor{
		Sigma:     0.046,
		Threshold: 0.6,
		Deviation: 1.0,
	}
	config.Default.Vtp = Transistor{
		Sigma:     0.046,
		Threshold: -0.6,
		Deviation: 1.0,
	}
	config.Default.Times = 5000

	as := assert.New(t)

	t.Run("001_Task.MkDir", func(t *testing.T) {
		task := NewTask()

		task.MkDir()
		expect := PathJoin(config.Default.SimulationDirectories.BaseDir,
			"Vtn0.6000-Sigma0.0460",
			"Vtp-0.6000-Sigma0.0460",
			"Times05000",
			"SEED00100")

		actual := task.SimulationDirectories.DstDir

		as.Equal(expect, actual, "they has equal")

		if _, err := os.Stat(expect); err != nil {
			as.Fail("cannot find", expect)
		} else {
			os.RemoveAll(expect)
		}
	})
}
