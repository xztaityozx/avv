package cmd

import "fmt"

type Task struct {
	SimulationFiles       SimulationFiles
	Vtn                   Transistor
	Vtp                   Transistor
	ParallelConfig        ParallelConfig
	AutoRemove            bool
	SimulationDirectories SimulationDirectories
	PlotPoint             PlotPoint
	SEED                  int
	Times                 int
}

func NewTask() Task {
	return config.Default
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
