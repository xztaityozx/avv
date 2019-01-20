package cmd

import "testing"

func TestNewTask(t *testing.T) {

	TU.Init()

	expect := Task{
		ParallelConfig: ParallelConfig{
			CountUp:  10,
			WaveView: 20,
			HSPICE:   30,
		},
		PlotPoint: PlotPoint{
			SignalNames: []string{"A", "B", "C"},
			Stop:        17.5,
			Step:        7.5,
			Start:       2.5,
		},
		SEED: SEED{
			Start: 1,
			End:   2000,
		},
		TaskNameFileName: "TaskName",
		SrcDir:           "Src",
		DstDir:           "Dst",
		Vtp: Transistor{
			Sigma:     0.7,
			Deviation: 1.0,
			Threshold: 0.2,
		},
		Vtn: Transistor{
			Threshold: 0.3,
			Deviation: 1.2,
			Sigma:     0.3,
		},
	}

	config.ParallelConfig = ParallelConfig{
		CountUp:  10,
		WaveView: 20,
		HSPICE:   30,
	}

	actual := NewTask("Dst", "Src", "TaskName", 1, 2000, 0.3, 0.3,
		1.2, 0.2, 0.7, 1.0, 2.5, 7.5, 17.5, []string{"A", "B", "C"})

	TU.Assert(actual.Compare(expect), t, "\n", actual, "\nis not\n", expect)
}
