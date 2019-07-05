package parameters

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/xztaityozx/avv/fileutils"
	"golang.org/x/xerrors"
)

type Templates struct {
	SPIScript string
}

type Files struct {
	AddFile       string
	SPIScript     string
	ACEScript     string
	ResultsXML    string
	ResultsMapXML string
	Directories   Directories
	ResultFile    string
	TaskFile      string
}

type Directories struct {
	DstDir     string
	NetListDir string
	BaseDir    string
	SearchDir  string
	ResultDir  string
}

func (d Directories) MakeDirectories() error {
	check := func(p string) error {
		if err := fileutils.TryMakeDirAll(p); err != nil {
			return xerrors.Errorf("Failed make directories : %w", err)
		}

		return nil
	}

	for _, v := range []string{d.DstDir, d.ResultDir} {
		err := check(v)
		if err != nil {
			return xerrors.Errorf("can not find directory: %s", v)
		}
	}

	return nil
}

func Generate(base, net string, parameters Parameters) (Files, error) {

	sha := parameters.Hash()
	d := Directories{
		BaseDir:    base,
		NetListDir: net,
		DstDir:     filepath.Join(base, "sim", sha, fmt.Sprint(parameters.Seed), "sim"),
		ResultDir:  filepath.Join(base, "sim", sha, fmt.Sprint(parameters.Seed), "res"),
	}

	for _, v := range []string{d.DstDir, d.ResultDir} {
		err := fileutils.TryMakeDirAll(v)
		if err != nil {
			return Files{}, err
		}
	}

	return Files{
		Directories:   d,
		AddFile:       filepath.Join(d.DstDir, "add"),
		SPIScript:     filepath.Join(net, parameters.HashWithSeed()+".spi"),
		ACEScript:     filepath.Join(d.DstDir, "ace"),
		ResultsXML:    filepath.Join(d.DstDir, "results.xml"),
		ResultsMapXML: filepath.Join(d.DstDir, "resultsMap.xml"),
		ResultFile:    filepath.Join(d.ResultDir, fmt.Sprintf("SEED%05d.csv", parameters.Seed)),
	}, nil

}

// GenerateSPIScript write spi script to path
// params:
//  - path: path to spi script
//  - search: path to search dir
//  - add: path to addfile
//  - parameters: Parameters struct
// returns:
//  - error: error
func (t Templates) GenerateSPIScript(path, search, add string, parameters Parameters) error {
	b, err := ioutil.ReadFile(t.SPIScript)
	if err != nil {
		return err
	}
	tmp := string(b)
	// .option search='/path/to/SearchDir'
	// .param vtn=AGAUSS(th,sig,dev) vtp=AGAUSS(th,sig,dev)
	// .include '/path/to/AddFile'
	// .include '/path/to/ModelFile'
	// .tran 10p 20n start=0 uic sweep monte=Times firstrun=1

	// make spi script from template string
	data := fmt.Sprintf(tmp,
		parameters.Vtn.Threshold, parameters.Vtn.Sigma, parameters.Vtn.Deviation,
		parameters.Vtp.Threshold, parameters.Vtp.Sigma, parameters.Vtp.Deviation,
		add,
		parameters.ModelFile,
		parameters.Sweeps)

	// write script
	return fileutils.WriteFile(path, data)
}
