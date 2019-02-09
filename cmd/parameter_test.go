package cmd

import (
	"context"
	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestParameter_InsertQuery(t *testing.T) {
	as := assert.New(t)
	expect := "insert or ignore into Parameter(VtnId, VtpId, Times, Signals) values (?,?,?,?)"
	actual := Parameter{}.InsertQuery()

	as.Equal(expect, actual)
}

func TestParameter_SelectQuery(t *testing.T) {
	as := assert.New(t)
	expect := "select * from Parameter where VtnId = ? and VtpId = ? and Times = ? and Signals = ?"
	actual := Parameter{}.SelectQuery()

	as.Equal(expect, actual)
}
func TestParameter_Insert(t *testing.T) {
	as := assert.New(t)
	home, _ := homedir.Dir()
	p := PathJoin(home, "TestDir", "db")
	FU.TryMkDir(p)
	p = PathJoin(p, "database.db")

	os.Remove(p)

	r := NewRepositoryFromFile(p)

	para := Parameter{
		Signals: `test`,
		Times:   2000,
		VtpId:   1,
		VtnId:   2,
	}

	err := para.Insert(context.Background(), r)
	as.NoError(err)

	t.Run("Check Select", func(t *testing.T) {
		res, err := para.Select(context.Background(), r)
		para.ParamsId = 1
		as.NoError(err)
		as.Equal(para, res)
	})
}
