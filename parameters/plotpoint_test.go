package parameters

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestPlotPoint_GenerateACEScript(t *testing.T) {
	pp := PlotPoint{
		Start:   0.1111,
		Step:    0.2222,
		Stop:    0.3333,
		Signals: []string{"A", "B", "C"},
	}

	path := "/path/to/csv.csv"

	expect := []byte(fmt.Sprintf(`
set xml [ sx_open_wdf "resultsMap.xml" ]
sx_current_sim_file $xml
set www [ sx_signal "A B C" ]
sx_export_csv on
sx_export_range 0.1111 0.3333 0.2222
sx_export_data "%s" $www
`, path))

	home, _ := homedir.Dir()
	p := filepath.Join(home, "TestDir", "ace")
	err := pp.GenerateACEScript(p, path)

	as := assert.New(t)
	as.NoError(err)

	actual, err := ioutil.ReadFile(p)
	as.NoError(err)

	as.Equal(expect, actual)
}
