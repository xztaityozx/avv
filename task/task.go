package task

import (
	"github.com/xztaityozx/avv/parameters"
	"github.com/xztaityozx/avv/simulation"
)

// ITask Dispatcherで扱いたいstructに実装するInterface
type ITask interface {
	// Run 実際の処理を記述する
}

// Task タスクを定義するstructの基底クラス()
type Task struct {
	// SimulationFiles このタスクで扱うパラメータ情報のファイルです
	SimulationFiles simulation.SimulationFiles
	// Vtn Vtnのトランジスタ情報
	Vtn parameters.Transistor
	// Vtp Vtpのトランジスタ情報
	Vtp parameters.Transistor
	// TODO: AutoRemoveを実装しようね
	// AutoRemove タスク終了時にゴミを掃除します
	AutoRemove bool
	// SimulationDirectories このタスクで扱うディレクトリ情報です
	SimulationDirectories simulation.SimulationDirectories
	// PlotPoint
	PlotPoint parameters.PlotPoint
	// SEED このタスクのSEEDです
	SEED int
	// Sweeps このタスクのMCSの回数です
	Sweeps int
	// ResultCSV 結果を書き出したCSVへのパスです
	ResultCSV []string
	// Failure このタスクの実行結果から得られた不良数です
	Failure int64
}

