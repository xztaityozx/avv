package cmd

import (
	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRepositoryFromFile(t *testing.T) {
	as := assert.New(t)
	home, _ := homedir.Dir()
	p := PathJoin(home, "TestDir", "db", "database")
	actual := NewRepositoryFromFile(p)
	expect := Repository{Path: p}

	as.Equal(expect, actual)
}

func TestNewRepository(t *testing.T) {
	as := assert.New(t)
	home, _ := homedir.Dir()
	config.Default.Repository.Path = PathJoin(home, "TestDir", "db", "database")
	actual := NewRepository()
	expect := config.Default.Repository

	as.Equal(expect, actual)
}

func TestRepository_Connect(t *testing.T) {
	as := assert.New(t)
	home, _ := homedir.Dir()
	p := PathJoin(home, "TestDir", "db")
	FU.TryMkDir(p)

	p = PathJoin(p, "database.db")
	r := NewRepositoryFromFile(p)

	db, err := r.Connect()
	defer db.Db.Close()

	as.NoError(err)
	as.NotNil(db)
	//os.Remove(p)
}
