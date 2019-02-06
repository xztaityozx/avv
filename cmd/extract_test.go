package cmd

import (
	"fmt"
	"github.com/magiconair/properties/assert"
	"github.com/mitchellh/go-homedir"
	"os"
	"testing"
)

func TestExtractTask_GetExtractCommand(t *testing.T) {
	wd,_:=os.Getwd()
	home,_:=homedir.Dir()
	config.WaveView=WaveViewConfig{
		Command:PathJoin(wd,"..","test","wv.sh"),
	}

	task := Task{
		SimulationDirectories:SimulationDirectories{
			DstDir:PathJoin(home,"TestDir","Dst"),
		},
	}
	task.SimulationFiles.ACEScript="ACE"

	FU.TryMkDir(task.SimulationDirectories.DstDir)

	et := ExtractTask{Task:task}
	actual := et.GetExtractCommand()
	expect := fmt.Sprintf("cd %s && %s -k -ace_no_gui %s > wv.log",
		task.SimulationDirectories.DstDir,
		config.WaveView.Command,
		task.SimulationFiles.ACEScript)

	assert.Equal(t, actual, expect)
}
