package cmd

type Task struct {
	DstDir           string
	SrcDir           string
	TaskNameFileName string
	SEED             SEED
	PlotPoint        PlotPoint
	ParallelConfig   ParallelConfig
	Vtp              Transistor
	Vtn              Transistor
}

type SEED struct {
	Start int
	End   int
}

// Compare func for SEED struct
func (s SEED) Compare(t SEED) bool {
	return s.End == t.End && s.Start == t.Start
}

type Transistor struct {
	Threshold float64
	Sigma     float64
	Deviation float64
}

// Compare func for Transistor struct
func (t Transistor) Compare(s Transistor) bool {
	return t.Sigma == s.Sigma && t.Deviation == s.Deviation && t.Threshold == s.Threshold
}

type PlotPoint struct {
	Start       float64
	Step        float64
	Stop        float64
	SignalNames []string
}

// Compare func for PlotPoint struct
func (s PlotPoint) Compare(t PlotPoint) bool {
	if len(s.SignalNames) != len(t.SignalNames) {
		return false
	}

	for i, v := range s.SignalNames {
		if v != t.SignalNames[i] {
			return false
		}
	}

	return s.Start == t.Start && s.Step == t.Step && s.Stop == t.Stop
}

// NewTask: タスクを作る，引数多すぎてウケる
func NewTask(dst, src, fname string,
	sstart, send int,
	nth, nsigma, ndev, pth, psigma, pdev, pstart, pstep, pstop float64,
	sigName []string) Task {
	return Task{
		DstDir:           dst,
		SrcDir:           src,
		TaskNameFileName: fname,
		SEED: SEED{
			Start: sstart,
			End:   send,
		},
		PlotPoint: PlotPoint{
			Start:       pstart,
			Step:        pstep,
			Stop:        pstop,
			SignalNames: sigName,
		},
		Vtp: Transistor{
			Threshold: pth,
			Deviation: pdev,
			Sigma:     psigma,
		},
		Vtn: Transistor{
			Threshold: nth,
			Deviation: ndev,
			Sigma:     nsigma,
		},
		ParallelConfig: config.ParallelConfig,
	}
}

// Compare func for Task struct
func (t Task) Compare(s Task) bool {
	return t.DstDir == s.DstDir &&
		t.SEED.Compare(s.SEED) &&
		t.TaskNameFileName == s.TaskNameFileName &&
		t.ParallelConfig.Compare(s.ParallelConfig) &&
		t.SrcDir == s.SrcDir &&
		t.Vtp.Compare(s.Vtp) &&
		t.Vtn.Compare(s.Vtn) &&
		t.PlotPoint.Compare(s.PlotPoint)
}
