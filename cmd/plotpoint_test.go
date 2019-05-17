package cmd

import (
	"os"
	"testing"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
)

func TestAllPlotPoint(t *testing.T) {
	home, _ := homedir.Dir()
	as := assert.New(t)

	pp := PlotPoint{
		Stop:  17.5,
		Step:  7.5,
		Start: 2.5,
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

func TestFilter_ToAwkStatement(t *testing.T) {
	f := Filter{
		Status:     []string{">=1", "<=2", ">3"},
		SignalName: "A",
	}

	actual := f.ToAwkStatement(1)
	expect := "$1>=1&&$2<=2&&$3>3"

	assert.Equal(t, expect, actual)
}

func TestFilter_Compare(t *testing.T) {
	f := Filter{
		Status:     []string{">=1", "<=2", ">3"},
		SignalName: "A",
	}

	assert.True(t, f.Compare(f))
	assert.False(t, f.Compare(Filter{}))
}
