package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

type PlotPoint struct {
	Start   float64
	Step    float64
	Stop    float64
	Signals []string
}

// Signal Name and filter values for count up failure stage
type Filter struct {
	SignalName string
	Status     []string
}

// convert to json for ResultRecord.Signals
// returns: json-string, error
func (p PlotPoint) ToJson() (string, error) {
	out, err := json.Marshal(p)
	if err != nil {
		return "", err
	}

	return string(out), nil
}

// compare func fof Filter struct
func (f Filter) Compare(t Filter) bool {
	if len(f.Status) != len(t.Status) {
		return false
	}

	for i, v := range f.Status {
		if v != t.Status[i] {
			return false
		}
	}

	return f.SignalName == t.SignalName
}

// ToAwkStatement generate awk statement for count up
// returns: awk statement string
func (f Filter) ToAwkStatement(start int) string {
	var rt []string
	for i, v := range f.Status {
		rt = append(rt, fmt.Sprintf("$%d%s", i+start, v))
	}

	return strings.Join(rt, "&&")
}

// Get plot step
//returns: Number of Plot Steps
func (p PlotPoint) PlotSteps() int {
	return (int)((p.Stop-p.Start)/p.Step + 1.0)
}

// Make ACE script from PlotPoint struct
// returns: [file path] [error]
func (p PlotPoint) MkACEScript(dst string) (string, error) {
	rt := fmt.Sprintf(`set xml [ sx_open_wdf "resultsMap.xml" ]
sx_export_csv on
sx_export_range %.2fns %.2fns %.2fns`, p.Start, p.Stop, p.Step)

	for _, v := range p.Signals {
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
	if len(p.Signals) != len(t.Signals) {
		return false
	}

	for i, v := range p.Signals {
		if v != t.Signals[i] {
			return false
		}
	}

	return p.Start == t.Start && p.Step == t.Step && p.Stop == t.Stop
}
