package cmd

import (
	"context"
	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
	"os"
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
	template := `search='%s'
%.4f %.4f %.4f
%.4f %.4f %.4f
include='%s'
include='%s'
monte=%d
`
	config.Templates.SPIScript = PathJoin(home, "template", "spi")
	FU.TryMkDir(PathJoin(home, "template"))
	FU.WriteFile(config.Templates.SPIScript, template)

	wd, _ := os.Getwd()

	config.HSPICE = HSPICEConfig{
		Command: PathJoin(wd, "..", "test", "hspice.sh"),
		Option:  "",
	}

	st := SimulationTask{Task: task}
	res := st.Run(context.Background())

	assert.True(t, res.Status)
}

func TestSimulationTask_Run(t *testing.T) {
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
	template := `search='%s'
%.4f %.4f %.4f
%.4f %.4f %.4f
include='%s'
include='%s'
monte=%d
`
	config.Templates.SPIScript = PathJoin(home, "template", "spi")
	FU.TryMkDir(PathJoin(home, "template"))
	FU.WriteFile(config.Templates.SPIScript, template)

	wd, _ := os.Getwd()

	config.HSPICE = HSPICEConfig{
		Command: PathJoin(wd, "..", "test", "hspice.sh"),
		Option:  "err",
	}

	st := SimulationTask{Task: task}
	res := st.Run(context.Background())

	assert.False(t, res.Status)
}