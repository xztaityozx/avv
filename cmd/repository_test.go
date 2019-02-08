package cmd

import (
	"context"
	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewRepositoryFromFile(t *testing.T) {
	as := assert.New(t)
	home, _ :=homedir.Dir()
	p := PathJoin(home, "TestDir","db","database")
	actual := NewRepositoryFromFile(p)
	expect := Repository{Path:p}

	as.Equal(expect,actual)
}

func TestNewRepository(t *testing.T) {
	as := assert.New(t)
	home, _ :=homedir.Dir()
	config.Default.Repository.Path = PathJoin(home, "TestDir","db","database")
	actual := NewRepository()
	expect := config.Default.Repository

	as.Equal(expect,actual)
}

func TestRepository_Connect(t *testing.T) {
	as := assert.New(t)
	home, _ :=homedir.Dir()
	p := PathJoin(home, "TestDir","db")
	FU.TryMkDir(p)

	p=PathJoin(p,"database.db")
	r := NewRepositoryFromFile(p)

	db, err := r.Connect()
	defer db.Db.Close()

	as.NoError(err)
	as.NotNil(db)
	os.Remove(p)
}


func TestRepository_InsertTransistors(t *testing.T) {
	var tr []Transistor
	for i:=0;i < 20; i++ {
		tr=append(tr, Transistor{
			Sigma:0.046,
			Deviation:1.0,
			Threshold:float64(i+1),
		})
	}

	as := assert.New(t)
	home, _ :=homedir.Dir()
	p := PathJoin(home, "TestDir","db")
	FU.TryMkDir(p)

	p=PathJoin(p,"database.db")

	os.Remove(p)

	r := NewRepositoryFromFile(p)

	err := r.InsertTransistors(context.Background(), tr...)
	as.NoError(err)

	t.Run("CheckID", func(t *testing.T) {
		ids,err := r.SelectTransistorIds(context.Background(), tr...)
		as.NoError(err)

		for i,v := range ids {
			as.Equal(int64(i+1),v)
		}
	})

}

