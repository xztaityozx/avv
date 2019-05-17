package cmd

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTransistor_Compare(t *testing.T) {
	t1 := Transistor{
		Sigma:        0.046,
		Threshold:    1,
		Deviation:    2,
		TransistorId: 3,
	}
	t2 := Transistor{
		Sigma:        0.046,
		Threshold:    4,
		TransistorId: 5,
		Deviation:    6,
	}

	as := assert.New(t)
	as.True(t1.Compare(t1))
	as.False(t1.Compare(t2))
}
