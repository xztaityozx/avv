package cmd

import (
	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestAllPlotPoint(t *testing.T) {
	home, _ := homedir.Dir()
	as := assert.New(t)

	pp := PlotPoint{
		SignalNames: []string{"A", "B", "C"},
		Stop:        17.5,
		Step:        7.5,
		Start:       2.5,
	}

	t.Run("001_MkACEScript", func(t *testing.T) {
		base := PathJoin(home, "TestDir")
		FU.TryMkDir(base)

		p, err := pp.MkACEScript(base)

		if err != nil {
			as.Fail("error has occur", err)
		}

		as.Equal(PathJoin(base, "extract.ace"), p)

		if _, err := os.Stat(p); err != nil {
			as.Fail(p, "could not found")
		}

	})
}
