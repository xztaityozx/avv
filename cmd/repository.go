package cmd

import (
	"context"
	"database/sql"
	"github.com/go-gorp/gorp"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"time"
)

// New Repository struct
func NewRepository() Repository {
	return config.Default.Repository
}

func NewRepositoryFromFile(p string) Repository {
	return Repository{Path: p}
}

func (r Repository) Connect() (dbMap *gorp.DbMap, err error) {
	if _, err := os.Stat(r.Path); err == nil {
		log.WithField("at", "Repository.CreateDB").Info("already exists database: ", r.Path)
	} else {
		_, err = os.Create(r.Path)
		if err != nil {
			return nil, err
		}

		log.WithField("at", "Repository.CreateDB").Info("create database ", r.Path)
	}

	db, err := sql.Open("sqlite3", r.Path)
	if err != nil {
		return nil, err
	}

	dbMap = &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	{
		table := dbMap.AddTable(Result{}).SetKeys(true, "Id")
		table.ColMap("ParamsId").SetNotNull(true)
		table.ColMap("Seed").SetNotNull(true)
		table.ColMap("Failure").SetNotNull(true)
	}
	{
		table := dbMap.AddTable(Parameter{}).SetKeys(true, "ParamsId")
		table.ColMap("VtnId").SetNotNull(true)
		table.ColMap("VtpId").SetNotNull(true)
		table.ColMap("Times").SetNotNull(true)
		table.ColMap("Signals").SetNotNull(true)
		table.SetUniqueTogether("VtnId", "VtpId", "Times", "Signals")
	}
	{
		table := dbMap.AddTable(Transistor{}).SetKeys(true, "TransistorId")
		table.ColMap("Threshold").SetNotNull(true)
		table.ColMap("Deviation").SetNotNull(true)
		table.ColMap("Sigma").SetNotNull(true)
		table.SetUniqueTogether("Threshold", "Deviation", "Sigma")
	}

	err = dbMap.CreateTablesIfNotExists()

	return
}

// insert or ignore some Transistor struct
// returns: error
func (r Repository) InsertTransistors(ctx context.Context, items ...Transistor) error {
	db, err := r.Connect()
	defer db.Db.Close()
	if err != nil {
		return err
	}

	stmt, err := db.Prepare(Transistor{}.InsertQuery()) // insert or ignore into Transistor(Deviation, Threshold, Sigma) values (?,?,?)
	for _, v := range items {
		// bind and execute query
		_, err := stmt.ExecContext(ctx, v.Deviation, v.Threshold, v.Sigma)
		if err != nil {
			return err
		}
	}

	return nil
}

// get Transistor's ids
// return: ids, error
func (r Repository) SelectTransistorIds(ctx context.Context, items ...Transistor) ([]int64,error) {
	var rt []int64

	db, err := r.Connect()
	defer db.Db.Close()
	if err != nil {
		return nil,err
	}

	for _, v := range items{
		// select TransistorId from Transistor where Deviation = ? and Threshold = ? and Sigma = ?
		id,err :=db.SelectInt(v.SelectIdQuery(), v.Deviation, v.Threshold, v.Sigma)
		if err != nil {
			return nil, err
		}

		rt = append(rt, id)
	}

	return rt,nil
}

func (r Repository) InsertParameters(ctx context.Context, items ...Parameter) error {

	return nil
}

type (
	Result struct {
		Id       int
		ParamsId int64
		Seed     int64
		Failure  int64
		Date     time.Time
	}

	Parameter struct {
		ParamsId int64
		VtnId    int64
		VtpId    int64
		Times    int64
		Signals  string
	}
	Repository struct {
		Path string
	}
)
