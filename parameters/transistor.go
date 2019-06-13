package parameters

import (
	"fmt"
)

type Transistor struct {
	TransistorId int64
	Threshold    float64
	Sigma        float64
	Deviation    float64
}

// Compare func for Transistor struct
func (t Transistor) Compare(s Transistor) bool {
	return t.Sigma == s.Sigma && t.Deviation == s.Deviation && t.Threshold == s.Threshold && t.TransistorId == s.TransistorId
}

func (t Transistor) StringPrefix(prefix string) string {
	return fmt.Sprintf("%s%.4f-Sigma%.4f", prefix, t.Threshold, t.Sigma)
}
