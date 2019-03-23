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

func TestPlotPoint_GetAwkScript(t *testing.T) {
	pp := PlotPoint{
		Filters: []Filter{
			{SignalName: "A", Status: []string{">=1", "<=6", ">1"}},
			{SignalName: "B", Status: []string{">=2", "<=7", ">2"}},
			{SignalName: "C", Status: []string{">=3", "<=8", ">3"}},
			{SignalName: "D", Status: []string{">=4", "<=9", ">4"}},
			{SignalName: "E", Status: []string{">=5", "<=0", ">5"}},
		},
		Start: 2.5,
		Step:  7.5,
		Stop:  17.5,
	}

	actual := pp.GetAwkScript()
	expect := "BEGIN{sum=0}$1>=1&&$2<=6&&$3>1&&$4>=2&&$5<=7&&$6>2&&$7>=3&&$8<=8&&$9>3&&$10>=4&&$11<=9&&$12>4&&$13>=5&&$14<=0&&$15>5{sum++}END{print sum}"

	assert.Equal(t, expect, actual)
}
