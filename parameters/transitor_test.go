package parameters

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTransistor_String(t *testing.T) {
	s := Transistor{
		name: "Vtn",
		Threshold:0.1,
		Sigma:0.2,
		Deviation:0.3,
	}

	actual := s.String()
	expect := fmt.Sprintf("Vtn%.4f-Deviation%.4f-Sigma%.4f", 0.1,0.3,0.2)

	assert.Equal(t, expect, actual)
}
