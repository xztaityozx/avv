package cmd

import (
	"fmt"
	"os"
)

type Task struct {
	DstDir           string
	SrcDir           string
	TaskNameFileName string
	AutoRemove       bool
	SEED             SEED
	PlotPoint        PlotPoint
	ParallelConfig   ParallelConfig
	Vtp              Transistor
	Vtn              Transistor
	Files            []SimulationFiles
	BaseDir          string
}

type SEED struct {
	Start int
	End   int
	Paths []string
}

// Compare func for SEED struct
func (s SEED) Compare(t SEED) bool {
	return s.End == t.End && s.Start == t.Start
}

// NewTask: タスクを作る，引数多すぎてウケる
func NewTask(dst, src, fname string,
	sstart, send int,
	nth, nsigma, ndev, pth, psigma, pdev float64,
	ar bool) Task {
	return Task{
		DstDir:           dst,
		SrcDir:           src,
		TaskNameFileName: fname,
		SEED: SEED{
			Start: sstart,
			End:   send,
		},
		PlotPoint: config.PlotPoint,
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
		AutoRemove:     ar,
	}
}

// Compare func for Task struct
func (t Task) Compare(s Task) bool {

	if len(t.Files) != len(s.Files) {
		return false
	}

	for i, v := range t.Files {
		if !v.Compare(s.Files[i]) {
			return false
		}
	}

	return t.DstDir == s.DstDir &&
		t.SEED.Compare(s.SEED) &&
		t.TaskNameFileName == s.TaskNameFileName &&
		t.ParallelConfig.Compare(s.ParallelConfig) &&
		t.SrcDir == s.SrcDir &&
		t.Vtp.Compare(s.Vtp) &&
		t.Vtn.Compare(s.Vtn) &&
		t.PlotPoint.Compare(s.PlotPoint) &&
		t.AutoRemove == s.AutoRemove

}

func (t *Task) MkDirs() error {
	log.Info("Task.MkDirs")
	log.Info("Dst: ", t.DstDir)

	base := PathJoin(t.DstDir, "RangeSEED", fmt.Sprintf("Vtn%.4f-Sigma%.4f", t.Vtn.Threshold, t.Vtn.Sigma),
		fmt.Sprintf("Vtp%.4f-Sigma%.4f", t.Vtp.Threshold, t.Vtp.Sigma))

	for s := t.SEED.Start; s <= t.SEED.End; s++ {
		p := PathJoin(base, fmt.Sprintf("SEED%05d.csv", s))
		FU.TryMkDir(p)
		t.SEED.Paths = append(t.SEED.Paths, p)
	}
	t.BaseDir = base

	_, err := os.Stat(base)

	return err
}

func (t *Task) MkFiles() error {
	log.Info("Task.MkFiles")

	acePath, err := t.PlotPoint.MkACEScript(t.BaseDir)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
