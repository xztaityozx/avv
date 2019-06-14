package simulation

import (
	"fmt"
	"github.com/xztaityozx/avv/fileutils"
	"strings"
)

// AddFile Addfileを表すstruct
type AddFile struct {
	// VddVoltage Vddの電圧
	VddVoltage float64
	// GndVoltage Gndの電圧
	GndVoltage float64
	// ICCommand
	ICCommand string
	// オプション
	Options []string
	// SEED Seed値
	SEED int
}

// GenerateAddFile is generate addfile for hspice simulation
// params:
//  - path: path to addfile
func (a AddFile) GenerateAddFile(path string) error {
	data := fmt.Sprintf(`VDD VDD! 0 %.1fV
VGND GND! 0 %.1fV
%s
%s
.option SEED=%d`, a.VddVoltage, a.GndVoltage, a.ICCommand, strings.Join(a.Options, "\n"), a.SEED)

	// Try make Directory
	return fileutils.WriteFile(path, data)
}
