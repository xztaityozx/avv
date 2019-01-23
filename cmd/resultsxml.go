package cmd

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

const (
	XMLmonteCarlo   = "MONTE_CARLO"
	XMLanalysisName = "analysisName"
	XMLfilename     = "filename"
	XMLformat       = "format"
	XMLgiString     = "giString"
	XMLname         = "name"
	XMLparams       = "params"
	XMLpsf          = "psf"
	XMLresultFiles  = "resultFiles"
	XMLresultType   = "resultType"
	XMLsaPair       = "saPair"
	XMLsaResultFile = "saResultFile"
	XMLsaResults    = "saResults"
	XMLsaSweepFile  = "saSweepFile"
	XMLstatistical  = "statistical"
	XMLsweepFiles   = "sweepFiles"
	XMLtran         = "tran"
	XMLvalue        = "value"
)

type Attribute struct {
	XMLName xml.Name `xml:"attribute"`
	Type    string   `xml:"type,attr"`
	Value   string   `xml:"value,attr"`
	Name    string   `xml:"name,attr"`
}

type Object struct {
	XMLName     xml.Name     `xml:"object"`
	Version     string       `xml:"version,attr"`
	Type        string       `xml:"type,attr"`
	Name        string       `xml:"name,attr"`
	Attributes  []Attribute  `xml:"Attribute"`
	Collections []Collection `xml:"collection"`
}

type Collection struct {
	XMLName xml.Name `xml:"collection"`
	Name    string   `xml:"name,attr"`
	Objects []Object `xml:"Object"`
}

type FileFormat struct {
	XMLName xml.Name `xml:"file-format"`
	Version string   `xml:"version,attr"`
	Name    string   `xml:"name,attr"`
	Objects []Object `xml:"Object"`
}

type ResultsXML struct {
	Task Task
}

func (r ResultsXML) MakeResultsFilesCollection() Collection {

	mkIteration := func(t int) string {
		var rt []string

		for i := 1; i <= t; i++ {
			rt = append(rt, fmt.Sprintf("%d.0", i))
		}

		return strings.Join(rt, " ")
	}

	rt := Collection{
		Name: XMLresultFiles,
		Objects: []Object{
			{
				Version: "1",
				Type:    XMLsaResultFile,
				Name:    XMLresultFiles,
				Attributes: []Attribute{
					{Type: XMLgiString, Value: XMLtran, Name: XMLanalysisName},
					{Type: XMLgiString, Value: "", Name: XMLfilename},
					{Type: XMLgiString, Value: XMLpsf, Name: XMLformat},
					{Type: XMLgiString, Value: mkIteration(r.Task.Times), Name: "iterations"},
					{Type: XMLgiString, Value: XMLtran, Name: XMLresultType},
				},
				Collections: []Collection{r.MakeSweepFilesCollections()},
			},
			{
				Version: "1",
				Type:    XMLsaResultFile,
				Name:    XMLresultFiles,
				Attributes: []Attribute{
					{Type: XMLgiString, Value: XMLstatistical, Name: XMLanalysisName},
					{Type: XMLgiString, Value: "hspice.mc0", Name: XMLfilename},
					{Type: XMLgiString, Value: XMLpsf, Name: XMLformat},
					{Type: XMLgiString, Value: XMLstatistical, Name: XMLresultType},
				},
			},
			{
				Version: "1",
				Type:    XMLsaResultFile,
				Name:    XMLresultFiles,
				Attributes: []Attribute{
					{Type: XMLgiString, Value: "scalarData", Name: XMLanalysisName},
					{Type: XMLgiString, Value: "scalar.dat", Name: XMLfilename},
					{Type: XMLgiString, Value: "table", Name: XMLformat},
					{Type: XMLgiString, Value: XMLstatistical, Name: XMLresultType},
				},
			},
			{
				Version: "1",
				Type:    XMLsaResultFile,
				Name:    XMLresultFiles,
				Attributes: []Attribute{
					{Type: XMLgiString, Value: "", Name: XMLanalysisName},
					{Type: XMLgiString, Value: "designVariables.wdf", Name: XMLfilename},
					{Type: XMLgiString, Value: "wdf", Name: XMLformat},
					{Type: XMLgiString, Value: "variables", Name: XMLresultType},
				},
			},
		},
	}

	return rt
}

func (r ResultsXML) MakeSweepFileObjects() []Object {
	var rt []Object
	for i := 1; i <= r.Task.Times; i++ {
		obj := Object{
			Version:    "1",
			Type:       XMLsaSweepFile,
			Name:       XMLsweepFiles,
			Attributes: []Attribute{{Type: XMLgiString, Value: fmt.Sprintf("hspice.tr0@%d", i), Name: XMLfilename}},
			Collections: []Collection{
				{
					Name: XMLparams,
					Objects: []Object{
						{
							Version: "1",
							Type:    XMLsaPair,
							Name:    XMLparams,
							Attributes: []Attribute{
								{Type: XMLgiString, Value: XMLmonteCarlo, Name: XMLname},
								{Type: XMLgiString, Value: fmt.Sprintf("%d.0", i), Name: XMLvalue},
							},
						},
					},
				},
			},
		}
		rt = append(rt, obj)
	}
	return rt
}

func (r ResultsXML) MakeSweepFilesCollections() Collection {

	rt := Collection{
		Name:    XMLsweepFiles,
		Objects: r.MakeSweepFileObjects(),
	}

	return rt
}

func (t *Task) NewResultsXml() (string, error) {
	netlist := t.SimulationDirectories.NetListDir
	dst := t.SimulationDirectories.DstDir

	if _, err := os.Stat(netlist); err != nil {
		return "", errors.New(fmt.Sprint("can not found ", netlist, " dir (NewResultsXml)"))
	}

	rx := ResultsXML{Task: *t}

	path := PathJoin(dst, "results.xml")

	data := FileFormat{
		Version: "1.0",
		Name:    XMLsaResults,
		Objects: []Object{
			{
				Version: "1",
				Type:    XMLsaResults,
				Name:    XMLsaResults,
				Attributes: []Attribute{
					{Type: XMLgiString, Value: "resultsMap.xml", Name: "name"},
					{Type: XMLgiString, Value: netlist, Name: "netlistDir"},
					{Type: XMLgiString, Value: ".", Name: "resultsDir"},
					{Type: XMLgiString, Value: time.Now().Format(time.ANSIC), Name: "runTime"},
					{Type: XMLgiString, Value: "HSPICE", Name: "simulator"},
					{Type: XMLgiString, Value: "", Name: "version"},
				},
				Collections: []Collection{rx.MakeResultsFilesCollection()},
			},
		},
	}

	b, err := xml.MarshalIndent(data, "", " ")
	if err != nil {
		return "", err
	}

	log.Info("NewResultsXml: Write to", path)
	if err := ioutil.WriteFile(path, b, 0644); err != nil {
		return "", err
	}

	return path, nil
}
