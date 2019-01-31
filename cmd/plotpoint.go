package cmd

import (
	"fmt"
	"io/ioutil"
)

type PlotPoint struct {
	Start   float64
	Step    float64
	Stop    float64
	Filters []Filter
}

// Signal Name and filter values for count up failure stage
type Filter struct {
	SignalName string
	Values     []float64
}

// compare func fof Filter struct
func (f Filter) Compare(t Filter) bool {
	if len(f.Values) != len(t.Values) {
		return false
	}

	for i, v := range f.Values {
		if v != t.Values[i] {
			return false
		}
	}

	return f.SignalName == t.SignalName
}

// Get plot step
//returns: Number of Plot Steps
func (p PlotPoint) Count() int {
	return (int)((p.Stop-p.Start)/p.Step + 1.0)
}

// Make ACE script from PlotPoint struct
// returns: [file path] [error]
func (p PlotPoint) MkACEScript(dst string) (string, error) {
	rt := fmt.Sprintf(`set xml [ sx_open_wdf "resultsMap.xml" ]
sx_export_csv on
sx_export_range %.2fns %.2fns %.2fns`, p.Start, p.Step, p.Stop)

	for _, v := range p.Filters {
		rt = fmt.Sprintf(`%s
set www [ sx_find_wave_in_file $xml %s ]
sx_export_data "%s.csv" $www`, rt, v.SignalName, v.SignalName)
	}

	path := PathJoin(dst, "extract.ace")
	return path, ioutil.WriteFile(path, []byte(rt), 0644)
}

// Compare func for PlotPoint struct
// return: compare result
func (p PlotPoint) Compare(t PlotPoint) bool {
	if len(p.Filters) != len(t.Filters) {
		return false
	}

	for i, v := range p.Filters {
		if v.Compare(t.Filters[i]) {
			return false
		}
	}

	return p.Start == t.Start && p.Step == t.Step && p.Stop == t.Stop
}
