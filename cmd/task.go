package cmd

import (
	"fmt"
)

type Task struct {
	SimulationFiles       SimulationFiles
	Vtn                   Transistor
	Vtp                   Transistor
	AutoRemove            bool
	SimulationDirectories SimulationDirectories
	PlotPoint             PlotPoint
	SEED                  int
	Times                 int
}

func NewTask() Task {
	return config.Default
}

func ReserveDir() string {
	return PathJoin(config.TaskDir, "reserve")
}

func DoneDir() string {
	return PathJoin(config.TaskDir, "done")
}

func FailedDir() string {
	return PathJoin(config.TaskDir, "failed")
}

func DustDir() string {
	return PathJoin(config.TaskDir, "dust")
}

// make DstDir
func (t *Task) MkDir() {
	log.Info("Task.Mkdir: Make Directories for simulations")
	dst := PathJoin(
		t.SimulationDirectories.BaseDir,
		fmt.Sprintf("Vtn%.4f-Sigma%.4f", t.Vtn.Threshold, t.Vtn.Sigma),
		fmt.Sprintf("Vtp%.4f-Sigma%.4f", t.Vtp.Threshold, t.Vtp.Sigma),
		fmt.Sprintf("Times%05d", t.Times),
		fmt.Sprintf("SEED%05d", t.SEED))

	FU.TryMkDir(dst)

	t.SimulationDirectories.DstDir = dst
}

func (t *Task) MkSimulationFiles() {

	// Make AddFile
	t.SimulationFiles.AddFile.Make(t.SimulationDirectories.BaseDir)

	// Make ACEScript
	var err error
	t.SimulationFiles.ACEScript, err = t.PlotPoint.MkACEScript(t.SimulationDirectories.DstDir)
	if err != nil {
		log.Fatal(err)
	}

	// Make SPIScript
	t.MakeSPIScript()

	// Make results.xml
	if path, err := t.MakeResultsXml(); err != nil {
		log.Fatal(err)
	} else {
		t.SimulationFiles.ResultsXML = path
	}

	// Make resultsMap.xml
	if path, err := t.MakeMapXml(); err != nil {
		log.Fatal(err)
	} else {
		t.SimulationFiles.ResultsMapXML = path
	}
}

// Make SPI script for simulation
func (t *Task) MakeSPIScript() {
	tmp := FU.Cat(config.Templates.SPIScript)
	// .option search='/path/to/SearchDir'
	// .param vtn=AGAUSS(th,dev,sig) vtp=AGAUSS(th,dev,sig)
	// .include '/path/to/AddFile'
	// .include '/path/to/ModelFile'
	// .tran 10p 20n start=0 uic sweep monte=Times firstrun=1

	// make spi script from template string
	data := fmt.Sprintf(tmp,
		t.SimulationDirectories.SearchDir,
		t.Vtn.Threshold, t.Vtn.Deviation, t.Vtn.Sigma,
		t.Vtp.Threshold, t.Vtp.Deviation, t.Vtp.Sigma,
		t.SimulationFiles.AddFile.Path,
		t.SimulationFiles.ModelFile,
		t.Times)

	// spi script path
	path := PathJoin(t.SimulationDirectories.NetListDir,
		fmt.Sprintf("Vtn%.4f-Sigma%.4f-Vtp%.4f-Sigma%.4f-Times%05d.spi",
			t.Vtn.Threshold, t.Vtn.Sigma, t.Vtp.Threshold, t.Vtp.Sigma, t.Times))

	// write script
	FU.WriteFile(path, data)
	// set script path to Task struct
	t.SimulationFiles.SPIScript = path
}
