package cmd

import (
	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResultsXML_MakeSweepFileObjects(t *testing.T) {
	as := assert.New(t)
	task := Task{
		Times:100,
	}

	home, _ := homedir.Dir()
	task.SimulationDirectories.DstDir=PathJoin(home,"Dst")
	FU.TryMkDir(PathJoin(home,"Dsr"))
	task.SimulationDirectories.NetListDir=PathJoin(home,"NetList")
	FU.TryMkDir(task.SimulationDirectories.NetListDir)

	path, err := task.NewResultsXml()
	if err != nil {
		as.Error(err)
	}

	expect := FU.Cat("../test/100.xml")
	actual := FU.Cat(path)

	as.Equal(expect,actual)

}