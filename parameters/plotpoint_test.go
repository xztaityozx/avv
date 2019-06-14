package parameters

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPlotPoint_GenerateACEScript(t *testing.T) {
	pp := PlotPoint{
		Start:0.1111,
		Step: 0.2222,
		Stop: 0.3333,
		Signals: []string{"A","B","C"},
	}

	path := "/path/to/csv.csv"

	expect := fmt.Sprintf(`
set xml [ sx_open_wdf "resultsMap.xml"]
sx_current_sim_file $xml
set www [ sx_signal "A B C" ]
sx_export_csv on
sx_export_range 0.11110n 0.33330n 0.22220n
sx_export_data "%s" $www
`, path)

	actual := pp.GenerateACEScript(path)

	as := assert.New(t)
	as.Equal(expect, actual)
}
