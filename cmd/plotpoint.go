package cmd

import (
	"fmt"
	"io/ioutil"
)

type PlotPoint struct {
	Start       float64
	Step        float64
	Stop        float64
	SignalNames []string
}

// Make ACE script from PlotPoint struct
// returns: [file path] [error]
func (p PlotPoint) MkACEScript(dst string) (string, error) {
	rt := fmt.Sprintf(`set xml [ sx_open_wdf "resultsMap.xml" ]
sx_export_csv on
sx_export_range %.2fns %.2fns %.2fns`, p.Start, p.Step, p.Stop)

	for _, v := range p.SignalNames {
		rt = fmt.Sprintf(`%s
set www [ sx_find_wave_in_file $xml %s ]
sx_export_data "%s.csv" $www`, rt, v, v)
	}

	path := PathJoin(dst, "extract.ace")
	return path, ioutil.WriteFile(path, []byte(rt), 0644)
}

// Compare func for PlotPoint struct
// return: compare result
func (p PlotPoint) Compare(t PlotPoint) bool {
	if len(p.SignalNames) != len(t.SignalNames) {
		return false
	}

	for i, v := range p.SignalNames {
		if v != t.SignalNames[i] {
			return false
		}
	}

	return p.Start == t.Start && p.Step == t.Step && p.Stop == t.Stop
}
