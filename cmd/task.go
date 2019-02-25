package cmd

import (
	"context"
	"fmt"
	"strings"
)

// ITask Dispatcherで扱いたいstructに実装するInterface
type ITask interface {
	// Run 実際の処理を記述する
	Run(context.Context) TaskResult
	// Self 基底クラス()のTaskを返す関数
	Self() Task
}

// Task タスクを定義するstructの基底クラス()
type Task struct {
	// SimulationFiles このタスクで扱うパラメータ情報のファイルです
	SimulationFiles SimulationFiles
	// Vtn Vtnのトランジスタ情報
	Vtn Transistor
	// Vtp Vtpのトランジスタ情報
	Vtp Transistor
	// TODO: AutoRemoveを実装しようね
	// AutoRemove タスク終了時にゴミを掃除します
	AutoRemove bool
	// SimulationDirectories このタスクで扱うディレクトリ情報です
	SimulationDirectories SimulationDirectories
	// PlotPoint
	PlotPoint PlotPoint
	// SEED このタスクのSEEDです
	SEED int
	// Times このタスクのMCSの回数です
	Times int
	// Stage このタスクがHSPICE、WaveView、CountUp、DBAccessのうちどれかを表します
	Stage Stage
	// ResultCSV 結果を書き出したCSVへのパスです
	ResultCSV []string
	// Repository このタスクの結果を書き込むDB情報です
	Repository Repository
	// TaskId このタスクのグループIDです
	TaskId int64
	// Failure このタスクの実行結果から得られた不良数です
	Failure int64
}

// Stage
type Stage string

const (
	HSPICE   Stage = "HSPICE"
	WaveView Stage = "WaveView"
	CountUp  Stage = "CountUp"
	DBAccess Stage = "DBAccess"
	Remove   Stage = "Remove"
)

// GetWrapper Task.StageをもとにITaskなstructを返します
// returns: ITask
func (t Task) GetWrapper() ITask {
	if t.Stage == HSPICE {
		return SimulationTask{
			Task: t,
		}
	} else if t.Stage == WaveView {
		return ExtractTask{
			Task: t,
		}
	} else if t.Stage == CountUp {
		return CountTask{
			Task: t,
		}
	} else if t.Stage == DBAccess {
		return DBAccessTask{
			Task: t,
		}
	} else if t.Stage == Remove {
		return RemoveTask{
			Task: t,
		}
	}
	return SimulationTask{}
}

// Run ITask 実装用ダミー
func (t Task) Run(ctx context.Context) TaskResult {
	return TaskResult{}
}

// String
func (t Task) String() string {
	return fmt.Sprint("Task: ", t.Times, "-", t.Vtn.StringPrefix("Vtn"), "-", t.Vtp.StringPrefix("Vtp"))
}

func (t Task) Self() Task {
	return t
}

// NewTask StageがHSPICEなTask structをconfigをもとに作ります
// returns: Task
func NewTask() Task {
	config.Default.Stage = HSPICE
	return config.Default
}

// ReserveDir 実行待ちなタスクファイルの置き場を返します
// returns: /path/to/ReserveDir
func ReserveDir() string {
	return PathJoin(config.TaskDir, "reserve")
}

// DoneDir 実行後のタスクファイルの置き場を返します
// returns: /path/to/DoneDir
func DoneDir() string {
	return PathJoin(config.TaskDir, "done")
}

// FailedDir 失敗したタスクファイルの置き場を返します
// returns: /path/to/FailedDir
func FailedDir() string {
	return PathJoin(config.TaskDir, "failed")
}

// DustDir
func DustDir() string {
	return PathJoin(config.TaskDir, "dust")
}

// MkDir シミュレーションに必要なディレクトリを作成します
func (t *Task) MkDir() {
	log.Info("Task.Mkdir: Make Directories for simulations")
	// HSPICEとかの出力先
	dst := PathJoin(
		t.SimulationDirectories.BaseDir,
		fmt.Sprintf("Vtn%.4f-Sigma%.4f", t.Vtn.Threshold, t.Vtn.Sigma),
		fmt.Sprintf("Vtp%.4f-Sigma%.4f", t.Vtp.Threshold, t.Vtp.Sigma),
		fmt.Sprintf("Times%05d", t.Times),
		fmt.Sprintf("SEED%05d", t.SEED))

	// 作成
	FU.TryMkDir(dst)

	// 結果を書き出すところ
	resultDir := PathJoin(
		t.SimulationDirectories.BaseDir,
		fmt.Sprintf("Vtn%.4f-Sigma%.4f", t.Vtn.Threshold, t.Vtn.Sigma),
		fmt.Sprintf("Vtp%.4f-Sigma%.4f", t.Vtp.Threshold, t.Vtp.Sigma),
		fmt.Sprintf("Times%05d", t.Times),
		fmt.Sprintf("TaskResult"))

	// 作成
	FU.TryMkDir(resultDir)

	// 設定
	t.SimulationDirectories.DstDir = dst
	t.SimulationDirectories.ResultDir = resultDir

	log.Info("dst: ",dst)
	log.Info("resultDir: ",resultDir)
}

// MkSimulationFiles シミュレーションに必要なパラメータファイルを生成します
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
	// .param vtn=AGAUSS(th,sig,dev) vtp=AGAUSS(th,sig,dev)
	// .include '/path/to/AddFile'
	// .include '/path/to/ModelFile'
	// .tran 10p 20n start=0 uic sweep monte=Times firstrun=1

	// make spi script from template string
	data := fmt.Sprintf(tmp,
		t.SimulationDirectories.SearchDir,
		t.Vtn.Threshold, t.Vtn.Sigma, t.Vtn.Deviation,
		t.Vtp.Threshold, t.Vtp.Sigma, t.Vtp.Deviation,
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
	log.WithField("at", "MakeSPIScript").Info("Write SPI script to ", path)
}

// Generate simulation command
// return: command string
func (t Task) GetSimulationCommand() string {
	var rt []string

	// append cd command
	rt = append(rt, fmt.Sprintf("cd %s &&", t.SimulationDirectories.DstDir))
	// append hspice command
	rt = append(rt, config.HSPICE.GetCommand(t.SimulationFiles.SPIScript))

	return strings.Join(rt, " ")
}

// Generate cd command
func (t Task) GetCdCommand() string {
	return fmt.Sprintf("cd %s &&", t.SimulationDirectories.DstDir)
}
