package cmd

import (
	"context"
	"database/sql"
	"github.com/go-gorp/gorp"
	_ "github.com/mattn/go-sqlite3"
	"os"
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
		//log.WithField("at", "Repository.CreateDB").Info("already exists database: ", r.Path)
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
		table := dbMap.AddTableWithName(ResultRecord{}, "Results").SetKeys(true, "Id")
		table.ColMap("TaskId").SetNotNull(true)
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
	{
		table := dbMap.AddTableWithName(TaskGroup{}, "Groups").SetKeys(true, "TaskId")
		table.ColMap("ParamsId").SetNotNull(true)
		table.ColMap("SeedStart").SetNotNull(true)
		table.ColMap("SeedEnd").SetNotNull(true)
		table.ColMap("Date").SetNotNull(true)
		table.SetUniqueTogether("ParamsId", "SeedStart", "SeedEnd", "Date")
	}

	err = dbMap.CreateTablesIfNotExists()

	return
}

type (
	IRecord interface {
		InsertQuery() string
		Insert(ctx context.Context, repository Repository) error
		SelectQuery() string
	}

	Repository struct {
		Path string
	}
)
