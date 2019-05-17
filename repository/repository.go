package repository

import (
	"github.com/xztaityozx/avv/cmd"
)

type record struct {
	Vtn    cmd.Transistor `bson:",inline"`
	Vtp    cmd.Transistor `bson:",inline"`
	Sweeps int            `bson:"sweeps"`
	Seed   int            `bson:"seed"`
	Values []float64      `bson:"values"`
	Time   float64        `bson:"time"`
}

type Repository struct {
	address   string
	tableName string
}

func NewRepository(address, tableName string) Repository {
	return Repository{
		address:   address,
		tableName: tableName,
	}
}
