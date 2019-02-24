package cmd

import (
	"os"
	"testing"

	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
)


func TestSimulationDirectories_Remove(t *testing.T) {
	as := assert.New(t)

	home, _ := homedir.Dir()

	d := SimulationDirectories{
		DstDir:     PathJoin(home, "TestDir", "Dst"),
		NetListDir: PathJoin(home, "TestDir", "NetList"),
		BaseDir:    PathJoin(home, "TestDir"),
		SearchDir:  PathJoin(home, "TestDir", "Search"),
		ResultDir:  PathJoin(home, "TestDir", "Result"),
	}

	FU.TryMkDir(d.DstDir)
	FU.TryMkDir(d.NetListDir)
	FU.TryMkDir(d.BaseDir)
	FU.TryMkDir(d.SearchDir)
	FU.TryMkDir(d.ResultDir)

	err := d.Remove()

	as.NoError(err)
	_, err = os.Stat(d.DstDir)
	as.Error(err)
	_, err = os.Stat(d.NetListDir)
	as.NoError(err)
	_, err = os.Stat(d.BaseDir)
	as.NoError(err)
	_, err = os.Stat(d.SearchDir)
	as.NoError(err)
	_, err = os.Stat(d.ResultDir)
	as.NoError(err)
}

func TestSimulationFiles_Remove(t *testing.T) {
	as := assert.New(t)
	home, _ := homedir.Dir()
	testDir := PathJoin(home, "TestDir")
	FU.TryMkDir(testDir)
	s := SimulationFiles{
		ACEScript: PathJoin(testDir, "ACE"),
		AddFile: AddFile{
			Path: PathJoin(testDir, "AddFile"),
		},
		ModelFile:     PathJoin(testDir, "Model"),
		ResultsMapXML: PathJoin(testDir, "ResultsMapXML"),
		ResultsXML:    PathJoin(testDir, "ResultsXML"),
		SPIScript:     PathJoin(testDir, "SPIScript"),
		Self:PathJoin(testDir,"self.json"),
	}

	os.Create(s.SPIScript)
	os.Create(s.ResultsXML)
	os.Create(s.ResultsMapXML)
	os.Create(s.ModelFile)
	os.Create(s.AddFile.Path)
	os.Create(s.ACEScript)
	os.Create(s.Self)

	err := s.Remove()
	as.NoError(err)

	_, err = os.Stat(s.SPIScript)
	as.NoError(err)
	_, err = os.Stat(s.ModelFile)
	as.NoError(err)
	_, err = os.Stat(s.AddFile.Path)
	as.Error(err)
}
