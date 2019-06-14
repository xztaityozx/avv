package task

import (
	"github.com/xztaityozx/avv/parameters"
	"github.com/xztaityozx/avv/simulation"
)

// Task タスクを定義するstructの基底クラス()
type Task struct {
	// SimulationFiles このタスクで扱うパラメータ情報のファイルです
	Files simulation.Files
	// Directories
	Directories simulation.Directories
	// Vtn Vtnのトランジスタ情報
	Vtn parameters.Transistor
	// Vtp Vtpのトランジスタ情報
	Vtp parameters.Transistor
	// AutoRemove タスク終了時にゴミを掃除します
	AutoRemove bool
	// PlotPoint
	PlotPoint parameters.PlotPoint
	// SEED このタスクのSEEDです
	SEED int
	// Sweeps このタスクのMCSの回数です
	Sweeps int
}

