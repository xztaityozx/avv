package cmd

import (
	"context"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestExtractTask_GetExtractCommand(t *testing.T) {
	wd, _ := os.Getwd()
	home, _ := homedir.Dir()
	config.WaveView = WaveViewConfig{
		Command: PathJoin(wd, "..", "test", "wv.sh"),
	}

	task := Task{
		SimulationDirectories: SimulationDirectories{
			DstDir: PathJoin(home, "TestDir", "Dst"),
		},
	}
	task.SimulationFiles.ACEScript = "ACE"

	FU.TryMkDir(task.SimulationDirectories.DstDir)

	et := ExtractTask{Task: task}
	actual := et.GetExtractCommand()
	expect := fmt.Sprintf("cd %s && %s -k -ace_no_gui %s &> wv.log",
		task.SimulationDirectories.DstDir,
		config.WaveView.Command,
		task.SimulationFiles.ACEScript)

	assert.Equal(t, actual, expect)
}

func TestExtractTask_Run(t *testing.T) {
	wd, _ := os.Getwd()
	home, _ := homedir.Dir()
	config.WaveView = WaveViewConfig{
		Command: PathJoin(wd, "..", "test", "wv.sh"),
	}
	task := Task{
		SimulationDirectories: SimulationDirectories{
			DstDir: PathJoin(home, "TestDir", "Dst"),
			NetListDir:PathJoin(home,"TestDir","NetList"),
		},
	}

	FU.TryMkDir(task.SimulationDirectories.DstDir)
	FU.TryMkDir(task.SimulationDirectories.NetListDir)
	et := ExtractTask{Task: task}

	res := et.Run(context.Background())
	as := assert.New(t)
	as.True(res.Status)
}
