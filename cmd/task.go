package cmd

type Task struct {
	SimulationFiles       SimulationFiles
	Vtn                   Transistor
	Vtp                   Transistor
	ParallelConfig        ParallelConfig
	AutoRemove            bool
	SimulationDirectories SimulationDirectories
	PlotPoint             PlotPoint
	SEED                  int
}

//func NewTask() Task {
//	return config.Default
//}
