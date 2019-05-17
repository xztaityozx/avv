package repository

import (
	"github.com/globalsign/mgo"
	"github.com/xztaityozx/avv/cmd"
	"github.com/xztaityozx/go-wvparser"
)

type record struct {
	Vtn    cmd.Transistor `bson:",inline"`
	Vtp    cmd.Transistor `bson:",inline"`
	Seed   int            `bson:"seed"`
	Signal string         `bson:"signal"`
	Values []float64      `bson:"values"`
	Time   float64        `bson:"time"`
}

type Repository struct {
	Address   string
	TableName string
}

func NewRepository(address, tableName string) Repository {
	return Repository{
		Address:   address,
		TableName: tableName,
	}
}

// TODO: implement
func (r Repository) BackUp(file string) error {

	return nil
}

//
func (r Repository) Insert(vtn, vtp cmd.Transistor, seed int, csv *wvparser.WVCsv) error {

	session, err := mgo.Dial(r.Address)
	if err != nil {
		return err
	}
	defer session.Close()

	if err := session.Ping(); err != nil {
		return err
	}

	db := session.DB("result")
	collection := db.C("records")

	selector := record{
		Vtn:    vtn,
		Vtp:    vtp,
		Signal: csv.Header.Signal,
		Seed:   seed,
	}

	data := record{
		Vtn:    vtn,
		Vtp:    vtp,
		Signal: csv.Header.Signal,
		Seed:   seed,
	}

	m := map[float64][]float64{}

	for _, data := range csv.Data {
		for time, f := range data.Values {
			if v, ok := m[time]; ok {
				m[time] = append(v, f)
			} else {
				m[time] = []float64{f}
			}
		}
	}

	for time, values := range m {
		selector.Time = time
		data.Time = time

		data.Values = values
		if _, err := collection.Upsert(selector, data); err != nil {
			return err
		}
	}

	return nil
}
