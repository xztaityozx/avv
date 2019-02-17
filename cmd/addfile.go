package cmd

import (
	"fmt"
	"strings"
)

type AddFile struct {
	VddVoltage float64
	GndVoltage float64
	ICCommand  string
	Options    []string
	SEED       int
	Path       string
}

func NewAddFile(s int) AddFile {
	rt := config.Default.SimulationFiles.AddFile
	rt.SEED = s

	return rt
}

// Compare func for AddFile
// return: compare result
func (a AddFile) Compare(b AddFile) bool {
	if len(a.Options) != len(b.Options) {
		return false
	}

	for i, v := range a.Options {
		if b.Options[i] != v {
			return false
		}
	}

	return a.SEED == b.SEED &&
		a.GndVoltage == b.GndVoltage &&
		a.ICCommand == b.ICCommand &&
		a.Path == b.Path &&
		a.VddVoltage == b.VddVoltage
}

func (a *AddFile) Make(base string) {
	data := fmt.Sprintf(`
VDD VDD! 0 %.1fV
VGND GND! 0 %.1fV
%s
%s
.option SEED=%d`, a.VddVoltage, a.GndVoltage, a.ICCommand, strings.Join(a.Options, "\n"), a.SEED)

	// Try make Directory
	p := PathJoin(base, "AddFiles")
	FU.TryMkDir(p)

	p = PathJoin(p, fmt.Sprintf("%d", a.SEED))
	FU.WriteFile(p, data)
	a.Path = p
}
