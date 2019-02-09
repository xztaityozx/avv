package cmd

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewResultRecord(t *testing.T) {
	as:=assert.New(t)
	actual := NewResultRecord(CountResult{
		Failure:20,
		Task:Task{TaskId:1, SEED:10},
	})

	as.Equal(20, actual.Failure)
	as.Equal(1, actual.TaskId)
	as.Equal(10,actual.Seed)
}
