package cmd

import (
	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestAddFile_Compare(t *testing.T) {
	as := assert.New(t)

	a := AddFile{
		SEED:       10,
		Options:    []string{"op1", "op2"},
		ICCommand:  "command",
		GndVoltage: 0.8,
		VddVoltage: 0.0,
		Path:       "/path/to",
	}

	as.True(a.Compare(a))
	b := a
	b.VddVoltage = 1.2

	as.False(a.Compare(b))
}

func TestAddFile_Make(t *testing.T) {
	as := assert.New(t)
	home, _ := homedir.Dir()
	base := PathJoin(home, "Base")
	a := AddFile{
		SEED:       10,
		Options:    []string{"op1", "op2"},
		ICCommand:  "command",
		GndVoltage: 1.2,
		VddVoltage: 0.8,
	}
	a.Make(base)

	if _, err := os.Stat(a.Path); err != nil {
		as.Fail("could not open addfile", err)
	}

	expect := `
VDD VDD! 0 0.8V
VGND GND! 0 1.2V
command
op1
op2
.option SEED=10`

	actual := FU.Cat(a.Path)

	as.Equal(expect, actual)
}

func TestNewAddFile(t *testing.T) {
	expect := config.Default.SimulationFiles.AddFile
	expect.SEED = 10
	actual := NewAddFile(10)

	assert.Equal(t, expect, actual)
}
