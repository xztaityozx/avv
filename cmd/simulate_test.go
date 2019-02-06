package cmd

import (
	"context"
	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSimulate_Run(t *testing.T) {
	home, _ := homedir.Dir()
	config.LogDir = PathJoin(home, "LogDir")
	task := Task{
		SimulationDirectories: SimulationDirectories{
			BaseDir:    PathJoin(home, "TestDir"),
			NetListDir: PathJoin(home, "TestDir", "NetList"),
			SearchDir:  PathJoin(home, "TestDir", "Search"),
		},
		Stage: HSPICE,
		SimulationFiles: SimulationFiles{
			AddFile: NewAddFile(1),
		},
	}

	FU.TryMkDir(config.LogDir)
	FU.TryMkDir(task.SimulationDirectories.BaseDir)
	FU.TryMkDir(task.SimulationDirectories.NetListDir)
	FU.TryMkDir(task.SimulationDirectories.SearchDir)

	config.HSPICE = HSPICEConfig{
		Command: "./test/hspice.sh",
		Option:  "",
	}

	st := SimulationTask{Task: task}
	res := st.Run(context.Background())

	assert.True(t, res.Status)
}
