package parameters

import (
	"fmt"
	"github.com/xztaityozx/avv/fileutils"
	"golang.org/x/xerrors"
	"io/ioutil"
	"path/filepath"
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
			return err
		}
	}

	return nil
}

func Generate(base, net, search string, parameters Parameters) Files {

	sha := parameters.Hash()
	d := Directories{
		BaseDir:    base,
		NetListDir: net,
		SearchDir:  search,
		DstDir:     filepath.Join(base, sha),
		ResultDir:  filepath.Join(base, sha, "Result"),
	}

	return Files{
		Directories:   d,
		AddFile:       filepath.Join(d.DstDir, "addfile"),
		SPIScript:     filepath.Join(net, sha+".spi"),
		ACEScript:     filepath.Join(d.DstDir, "ace"),
		ResultsXML:    filepath.Join(d.DstDir, "results.xml"),
		ResultsMapXML: filepath.Join(d.DstDir, "resultsMap.xml"),
		ResultFile:    filepath.Join(d.ResultDir, fmt.Sprintf("SEED%05d.csv", parameters.Seed)),
	}

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
		search,
		parameters.Vtn.Threshold, parameters.Vtn.Sigma, parameters.Vtn.Deviation,
		parameters.Vtp.Threshold, parameters.Vtp.Sigma, parameters.Vtp.Deviation,
		add,
		parameters.ModelFile,
		parameters.Sweeps)

	// write script
	return fileutils.WriteFile(path, data)
}
