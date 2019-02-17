package cmd

import (
	"github.com/mattn/go-pipeline"
	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestShapingCSV(t *testing.T) {
	as := assert.New(t)
	out, err := pipeline.Output(
		[]string{"seq", "200"},
		[]string{"shuf"},
		[]string{"xargs", "-n2"},
		[]string{"sed", "s/ /, /g"})
	as.NoError(err)
	home, _ := homedir.Dir()
	p := PathJoin(home, "TestDir", "test.csv")
	FU.TryMkDir(PathJoin(home, "TestDir"))
	err = ioutil.WriteFile(p, out, 0644)
	as.NoError(err)

	_, err = ShapingCSV(p, "signal", 4)
	as.NoError(err)
}
