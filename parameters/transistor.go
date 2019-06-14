package parameters

import (
	"fmt"
)

type Transistor struct {
	name string
	Threshold    float64
	Sigma        float64
	Deviation    float64
}

func (t Transistor) String() string {
	return fmt.Sprintf("%s%.4f-Deviation%.4f-Sigma%.4f", t.name, t.Threshold,t.Deviation, t.Sigma)
}
