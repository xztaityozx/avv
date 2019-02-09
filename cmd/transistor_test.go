package cmd

import (
	"context"
	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
	"os"
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

func TestTransistor_InsertQuery(t *testing.T) {
	as := assert.New(t)
	actual := Transistor{}.InsertQuery()
	expect := "insert or ignore into Transistor(Deviation, Threshold, Sigma) values (?,?,?)"
	as.Equal(expect, actual)
}

func TestTransistor_SelectQuery(t *testing.T) {
	as := assert.New(t)
	actual := Transistor{}.SelectQuery()
	expect := "select * from Transistor where Deviation = ? and Threshold = ? and Sigma = ?"
	as.Equal(expect, actual)
}

func TestTransistor_Insert(t *testing.T) {
	as := assert.New(t)
	home, _ := homedir.Dir()
	p := PathJoin(home, "TestDir", "db")
	FU.TryMkDir(p)
	p = PathJoin(p, "database.db")

	os.Remove(p)

	r := NewRepositoryFromFile(p)

	tr := Transistor{
		Deviation: 1.0,
		Threshold: 0.6,
		Sigma:     0.046,
	}

	err := tr.Insert(context.Background(), r)
	as.NoError(err)

	t.Run("Check_Select", func(t *testing.T) {
		res, err := tr.Select(context.Background(), r)
		tr.TransistorId = 1
		as.NoError(err)
		as.Equal(tr, res)
	})
}
