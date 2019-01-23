package cmd

import "fmt"

type Transistor struct {
	Threshold float64
	Sigma     float64
	Deviation float64
}

// Compare func for Transistor struct
func (t Transistor) Compare(s Transistor) bool {
	return t.Sigma == s.Sigma && t.Deviation == s.Deviation && t.Threshold == s.Threshold
}

func (t Transistor) ToString(prefix string) string {
	return fmt.Sprintf("prefix%.4f-Sigma%.4f", t.Threshold, t.Sigma)
}
