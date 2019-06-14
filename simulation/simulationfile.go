package simulation

import (
	"errors"
	"fmt"
	"github.com/xztaityozx/avv/fileutils"
	"github.com/xztaityozx/avv/parameters"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
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
	ModelFile     string
	Self          string
	HSpiceLogFile string
	directories Directories
}

// Generate is generate files for simulation
// returns:
//  - Files: Files struct that set AddFile, SPIScript, ACEScript, ResultsXML, ResultsMapXML, HSpiceLogFile
//  - error: error
func Generate(baseDir,netListDir,searchDir string, a AddFile, pp parameters.PlotPoint, xml ResultsXML, vtn, vtp parameters.Transistor, templates Templates, sweeps int) (Files, error) {
	rt := Files{}

	rt.directories = Directories{
		SearchDir:searchDir,
		BaseDir:baseDir,
		NetListDir:netListDir,

	}

	rt.HSpiceLogFile = filepath.Join(d.DstDir, "hspice.log")

	// AddFile
	{
		rt.AddFile = filepath.Join(d.DstDir, "AddFile")
		err := a.GenerateAddFile(rt.AddFile)
		if err != nil {
			return Files{}, err
		}
	}

	// ACE
	{
		out := filepath.Join(d.ResultDir, fmt.Sprintf("SEED%05d.csv", a.SEED))
		str := pp.GenerateACEScript(out)
		rt.ACEScript = filepath.Join(d.DstDir, "ace")
		err := fileutils.WriteFile(rt.ACEScript, str)
		if err != nil {
			return Files{}, err
		}
	}

	// ResultsXML
	{
		var err error
		rt.ResultsXML, rt.ResultsMapXML, err = xml.Generate()
		if err != nil {
			return Files{}, err
		}
	}

	// SPI
	{
		rt.SPIScript = filepath.Join(d.NetListDir, fmt.Sprintf("%s%s-Sweeps%05d.spi", vtn.String(), vtp.String(), sweeps))
		err := makeSPIScript(rt.SPIScript, templates.SPIScript, vtn, vtp, rt.AddFile, rt.ModelFile, d.SearchDir, sweeps)
		if err != nil {
			return Files{}, err
		}
	}

	return rt, nil
}

type Directories struct {
	DstDir     string
	NetListDir string
	BaseDir    string
	SearchDir  string
	ResultDir  string
}

func makeSPIScript(path, templateFile string, vtn, vtp parameters.Transistor, addFile, modelFile, searchDir string, sweeps int) error {

	tmp, err := ioutil.ReadFile(templateFile)
	if err != nil {
		return err
	}
	// .option search='/path/to/SearchDir'
	// .param vtn=AGAUSS(th,sig,dev) vtp=AGAUSS(th,sig,dev)
	// .include '/path/to/AddFile'
	// .include '/path/to/ModelFile'
	// .tran 10p 20n start=0 uic sweep monte=Times firstrun=1

	data := fmt.Sprintf(string(tmp),
		searchDir,
		vtn.Threshold, vtn.Sigma, vtn.Deviation,
		vtp.Threshold, vtp.Sigma, vtp.Deviation,
		addFile,
		modelFile,
		sweeps)

	return fileutils.WriteFile(path, data)
}

func (d Directories) makeDirectories() error {
	if err := fileutils.TryMakeDirAll(d.DstDir); err != nil {
		return err
	}

	if err := fileutils.TryMakeDirAll(d.ResultDir); err != nil {
		return err
	}

	return d.existsAll()
}

func (d Directories) existsAll() error {
	f := func(p string) bool {
		if _, err := os.Stat(p); err != nil {
			return false
		}
		return true
	}

	var ng []string
	if !f(d.DstDir) {
		ng = append(ng, d.DstDir)
	}

	if !f(d.ResultDir) {
		ng = append(ng, d.ResultDir)
	}

	if !f(d.SearchDir) {
		ng = append(ng, d.SearchDir)
	}

	if !f(d.NetListDir) {
		ng = append(ng, d.NetListDir)
	}

	if !f(d.BaseDir) {
		ng = append(ng, d.BaseDir)
	}

	if len(ng) != 0 {
		return errors.New(fmt.Sprintf("%s not found", strings.Join(ng, ",")))
	}
	return nil
}
