package parameters

import (
	"crypto/sha256"
	"fmt"
)

type Parameters struct {
	PlotPoint PlotPoint
	Seed      int
	Sweeps    int
	Vtn       Transistor
	Vtp       Transistor
	AddFile   AddFile
	ModelFile string
}

// GenerateParameters generate some Parameters struct from config
func GenerateParameters(config Config) []Parameters {
	var rt []Parameters
	base := config.Default.Parameters

	for i := config.Default.SEED.Start; i <= config.Default.SEED.End; i++ {
		base.Seed = i
		base.AddFile.SEED = i
		rt = append(rt, base)
	}

	return rt
}

// Hash generate hash string from parameters
func (parameters Parameters) Hash() string {
	return fmt.Sprint(sha256.Sum256([]byte(fmt.Sprintf("%s%s%s%010d%010d",
		parameters.PlotPoint.String(), parameters.Vtn.String(), parameters.Vtp.String(),
		parameters.Seed, parameters.Sweeps))))
}
