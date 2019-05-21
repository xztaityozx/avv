package repository

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/xztaityozx/go-wvparser"
)

type record struct {
	VtnDeviation float64   `bson:"vtn-dev"`
	VtnSigma     float64   `bson:"vtn-sigma"`
	VtnThreshold float64   `bson:"vtn-th"`
	VtpDeviation float64   `bson:"vtp-dev"`
	VtpSigma     float64   `bson:"vtp-sigma"`
	VtpThreshold float64   `bson:"vtp-th"`
	Seed         int       `bson:"seed"`
	Signal       string    `bson:"signal"`
	Values       []float64 `bson:"values"`
	Time         float64   `bson:"time"`
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
func (r Repository) Insert(vtn, vtp Transistor, seed int, csv *wvparser.WVCsv) error {

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

	selector := bson.M{
		"vtn-sigma": vtn.Sigma,
		"vtn-dev":vtn.Deviation,
		"vtn-th":vtn.Threshold,
		"vtp-sigma": vtp.Sigma,
		"vtp-dev":vtp.Deviation,
		"vtp-th":vtp.Threshold,
	}

	data := record{
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
