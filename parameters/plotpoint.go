package parameters

import (
	"fmt"
	"github.com/xztaityozx/avv/fileutils"
	"strings"
)

// PlotPoint is setting for extracting
type PlotPoint struct {
	Start   float64
	Step    float64
	Stop    float64
	Signals []string
}

// GenerateACEScript generate ace script for custom waveview
// params:
// 	- path string: path to csv file exported by custom waveview
func (p PlotPoint) GenerateACEScript(path, storePath string) error {
	str := fmt.Sprintf(`
set xml [ sx_open_wdf "resultsMap.xml"]
sx_current_sim_file $xml
set www [ sx_signal "%s" ]
sx_export_csv on
sx_export_range %.5fn %.5fn %.5fn
sx_export_data "%s" $www
`, strings.Join(p.Signals, " "), p.Start, p.Stop, p.Step, storePath)

	return fileutils.WriteFile(path, str)
}

func (p PlotPoint) String() string {
	return fmt.Sprintf("%.4f-%.4f-%.4f-%s", p.Start, p.Step, p.Stop, strings.Join(p.Signals, "-"))
}
