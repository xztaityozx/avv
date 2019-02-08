package cmd

import (
	"fmt"
)

type Transistor struct {
	TransistorId int
	Threshold float64
	Sigma     float64
	Deviation float64
}

// Compare func for Transistor struct
func (t Transistor) Compare(s Transistor) bool {
	return t.Sigma == s.Sigma && t.Deviation == s.Deviation && t.Threshold == s.Threshold
}

func (t Transistor) StringPrefix(prefix string) string {
	return fmt.Sprintf("%s%.4f-Sigma%.4f", prefix, t.Threshold, t.Sigma)
}
func (t Transistor) SelectIdQuery() string {
	return fmt.Sprintf("select TransistorId from Transistor where Deviation = ? and Threshold = ? and Sigma = ?")
}

func (t Transistor) InsertQuery() string{
	return "insert or ignore into Transistor(Deviation, Threshold, Sigma) values (?,?,?)"
}