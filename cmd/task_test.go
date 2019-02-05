package cmd

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestAllTask(t *testing.T) {

	home, _ := homedir.Dir()
	config.Default.SimulationDirectories.BaseDir = PathJoin(home, "Base")
	config.Default.SimulationDirectories.NetListDir = PathJoin(home, "NetList")

	FU.TryMkDir(config.Default.SimulationDirectories.BaseDir)
	FU.TryMkDir(config.Default.SimulationDirectories.NetListDir)

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
		}
	})
}

func TestTask_MakeSPIScript(t *testing.T) {
	template := `search='%s'
%.4f %.4f %.4f
%.4f %.4f %.4f
include='%s'
include='%s'
monte=%d
`

	task := config.Default
	home, _ := homedir.Dir()
	config.Templates.SPIScript = PathJoin(home, "template", "spi")
	FU.TryMkDir(PathJoin(home, "template"))
	FU.WriteFile(config.Templates.SPIScript, template)

	task.MakeSPIScript()
	read := FU.Cat(task.SimulationFiles.SPIScript)
	assert.Equal(t, fmt.Sprintf(template,
		task.SimulationDirectories.SearchDir,
		task.Vtn.Threshold, task.Vtn.Deviation, task.Vtn.Sigma,
		task.Vtp.Threshold, task.Vtp.Deviation, task.Vtp.Sigma,
		task.SimulationFiles.AddFile.Path,
		task.SimulationFiles.ModelFile,
		task.Times), read)
}
