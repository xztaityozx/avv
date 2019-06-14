package simulation

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
	sweeps     int
	netListDir string
}

func NewResultsXML(sweeps int, d Directories) ResultsXML {
	return ResultsXML{
		sweeps:     sweeps,
		netListDir: d.NetListDir,
	}
}

// Generate write xml files for simulations
func (r ResultsXML) Generate(xml, mapXml string) error {
	err := r.makeMapXml(xml)
	if err != nil {
		return err
	}

	return r.makeMapXml(mapXml)
}

func (r ResultsXML) makeResultsFilesCollection() Collection {

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
					{Type: XMLgiString, Value: mkIteration(r.sweeps), Name: "iterations"},
					{Type: XMLgiString, Value: XMLtran, Name: XMLresultType},
				},
				Collections: []Collection{r.makeSweepFilesCollections()},
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

func (r ResultsXML) makeSweepFileObjects() []Object {
	var rt []Object
	for i := 1; i <= r.sweeps; i++ {
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

func (r ResultsXML) makeSweepFilesCollections() Collection {

	rt := Collection{
		Name:    XMLsweepFiles,
		Objects: r.makeSweepFileObjects(),
	}

	return rt
}

func (r ResultsXML) makeResultsXml(path string) error {
	netList := r.netListDir

	if _, err := os.Stat(netList); err != nil {
		return errors.New(fmt.Sprint("can not found ", netList, " dir (makeResultsXml)"))
	}

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
					{Type: XMLgiString, Value: netList, Name: "netlistDir"},
					{Type: XMLgiString, Value: ".", Name: "resultsDir"},
					{Type: XMLgiString, Value: time.Now().Format(time.ANSIC), Name: "runTime"},
					{Type: XMLgiString, Value: "HSPICE", Name: "simulator"},
					{Type: XMLgiString, Value: "", Name: "version"},
				},
				Collections: []Collection{r.makeResultsFilesCollection()},
			},
		},
	}

	b, err := xml.MarshalIndent(data, "", " ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, b, 0644)

}

func (r ResultsXML) makeMapXml(path string) error {

	saResultsMap := "saResultsMap"

	data := FileFormat{
		Version: "1.0",
		Name:    saResultsMap,
		Objects: []Object{
			{
				Version: "1",
				Type:    saResultsMap,
				Name:    saResultsMap,
				Attributes: []Attribute{
					{Type: XMLgiString, Value: ".", Name: "masterResultsDir"},
					{Type: XMLgiString, Value: XMLmonteCarlo, Name: "monteCarlo"},
					{Type: XMLgiString, Value: "resultsMap.xml", Name: XMLname},
					{Type: XMLgiString, Value: ".", Name: "resultsMapDir"},
					{Type: XMLgiString, Value: "HSPICE", Name: "simulator"},
					{Type: XMLgiString, Value: time.Now().Format(time.ANSIC), Name: "timeStamp"},
				},
				Collections: []Collection{
					{
						Name: "resultsInfo",
						Objects: []Object{
							{
								Version: "1",
								Type:    "saResultsInfo",
								Name:    "resultsInfo",
								Attributes: []Attribute{
									{Type: XMLgiString, Value: r.netListDir, Name: "netlistDir"},
									{Type: XMLgiString, Value: ".", Name: "resultsDir"},
								},
							},
						},
					},
				},
			},
		},
	}

	b, err := xml.MarshalIndent(data, "", " ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, b, 0644)

}
