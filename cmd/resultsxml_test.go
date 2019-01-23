package cmd

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestResultsXML_MakeSweepFileObjects(t *testing.T) {
	as := assert.New(t)
	task := Task{
		Times: 100,
	}

	home, _ := homedir.Dir()
	home = PathJoin(home, "Testdir")
	task.SimulationDirectories.DstDir = PathJoin(home, "Dst")
	FU.TryMkDir(PathJoin(home, "Dst"))
	task.SimulationDirectories.NetListDir = "../netlist"
	wd, _ := os.Getwd()
	FU.TryChDir(task.SimulationDirectories.DstDir)
	FU.TryMkDir(task.SimulationDirectories.NetListDir)

	path, err := task.NewResultsXml()
	if err != nil {
		as.FailNow("", err)
	}

	FU.TryChDir(wd)

	expect := fmt.Sprintf(FU.Cat("../test/100.xml"), time.Now().Format(time.ANSIC))
	actual := FU.Cat(path)

	as.Equal(expect, actual)

}
