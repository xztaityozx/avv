package parameters

import (
	"crypto/sha256"
	"fmt"
	"path/filepath"
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
// params:
//  - seeds: range of seed
//  - config: config for avv
// returns:
//  - []Parameters: generated parameters
func GenerateParameters(seeds SEED, config Config) []Parameters {
	var rt []Parameters
	base := config.Default.Parameters

	for i := seeds.Start; i <= seeds.End; i++ {
		base.Seed = i
		base.AddFile.SEED = i
		rt = append(rt, base)
	}

	return rt
}

// Hash generate hash string from parameters without seed value
// returns:
//  - string: hash string
func (parameters Parameters) Hash() string {

	//head := sha256.Sum256([]byte(
	//	fmt.Sprintf("%s%s%s%010d", parameters.PlotPoint.String(), parameters.Vtn.String(), parameters.Vtp.String(), parameters.Sweeps)))
	//
	//return filepath.Join(fmt.Sprintf("%x", head))
	return fmt.Sprintf("%s%s%s%010d", parameters.PlotPoint.String(), parameters.Vtn.String(), parameters.Vtp.String(), parameters.Sweeps)
}

// HashWithSeed generate hash string from parameters with seed value
// returns:
//  - string: hash string
func (parameters Parameters) HashWithSeed() string {
	head := sha256.Sum256([]byte(fmt.Sprintf("%s%s%s%010d%010d",
		parameters.PlotPoint.String(), parameters.Vtn.String(), parameters.Vtp.String(), parameters.Sweeps, parameters.Seed)))

	return filepath.Join(fmt.Sprintf("%x", head))

}
