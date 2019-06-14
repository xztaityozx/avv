package parameters

import (
	"fmt"
	"strings"
)

// PlotPoint is setting for extracting
type PlotPoint struct {
	Start   float64
	Step    float64
	Stop    float64
	Signals []string
}

// GenerateACEScript is generate ace script for custom waveview
// params:
// 	- path string: path to csv file exported by custom waveview
func (p PlotPoint) GenerateACEScript(path string) string {
	return fmt.Sprintf(`
set xml [ sx_open_wdf "resultsMap.xml"]
sx_current_sim_file $xml
set www [ sx_signal "%s" ]
sx_export_csv on
sx_export_range %.5fn %.5fn %.5fn
sx_export_data "%s" $www
`, strings.Join(p.Signals, " "), p.Start, p.Stop, p.Step, path)
}
