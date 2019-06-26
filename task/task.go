package task

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/xztaityozx/avv/parameters"
	"golang.org/x/xerrors"
)

type Task struct {
	// SimulationFiles このタスクで扱うパラメータ情報のファイルです
	Files parameters.Files
	// AutoRemove タスク終了時にゴミを掃除します
	AutoRemove bool
	// パラメータ
	parameters.Parameters
}

// Generate is generate some Task struct
func Generate(start, end int, config parameters.Config) ([]Task, error) {
	var rt []Task

	stat := func(p string) error {
		if _, err := os.Stat(p); err != nil {
			return errors.New(fmt.Sprintf("%s not found", p))
		}
		return nil
	}

	// check dirs
	if err := stat(config.Default.BaseDir); err != nil {
		return nil, err
	}

	if err := stat(config.Default.NetListDir); err != nil {
		return nil, err
	}

	if err := stat(config.Default.SearchDir); err != nil {
		return nil, err
	}

	// generate parameters
	for _, v := range parameters.GenerateParameters(parameters.SEED{start, end}, config) {
		f, err := parameters.Generate(config.Default.BaseDir, config.Default.NetListDir, config.Default.SearchDir, v)
		if err != nil {
			return rt, err
		}

		rt = append(rt, Task{
			Files:      f,
			AutoRemove: config.AutoRemove,
			Parameters: v,
		})
	}

	return rt, nil
}

// Unmarshal is unmarshal task file
// params:
//  - path: path to task file
// returns:
//  - Task : Task struct
//  - error: error
func Unmarshal(path string) (Task, error) {
	var rt Task
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return Task{}, err
	}

	err = json.Unmarshal(b, &rt)
	return rt, err
}

// MakeFiles make files and directories for simulation
func (t *Task) MakeFiles(tmp parameters.Templates) error {
	// Generate Directories
	err := t.Files.Directories.MakeDirectories()
	if err != nil {
		return err
	}

	// AddFile
	err = t.Parameters.AddFile.GenerateAddFile(t.Files.AddFile)
	if err != nil {
		return xerrors.Errorf("failed make AddFile: %s", err)
	}

	// SPIScript
	err = tmp.GenerateSPIScript(t.Files.SPIScript, t.Files.Directories.SearchDir, t.Files.AddFile, t.Parameters)
	if err != nil {
		return xerrors.Errorf("failed make SPIScript: %s", err)
	}

	// ACEScript
	err = t.PlotPoint.GenerateACEScript(t.Files.ACEScript, t.Files.ResultFile)
	if err != nil {
		return xerrors.Errorf("failed make ACEScript: %s", err)
	}

	//XML
	{
		xml := parameters.NewResultsXML(t.Sweeps, t.Files.Directories)
		err := xml.Generate(t.Files.ResultsXML, t.Files.ResultsMapXML)
		if err != nil {
			return err
		}
	}

	return nil
}
