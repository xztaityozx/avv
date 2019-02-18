package cmd

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"github.com/go-gorp/gorp"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"io"
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
func (r Repository) DBBackUp() error {

	if !config.AutoDBBackUp {
		logrus.Info("このタスクではDB(" + r.Path + ")にアクセスします。バックアップを作成しますか？")
		fmt.Printf("[y] はい [n] いいえ\n>>> ")
		s := bufio.NewScanner(os.Stdin)
		s.Scan()
		ans := s.Text()

		if ans != "y" {
			return nil
		}
	}

	dst := PathJoin(config.BackUpDir, r.Path+time.Now().Format("2006-01-02-15-04-05"))

	dfp, err := os.OpenFile(dst,os.O_CREATE|os.O_WRONLY,0644)
	defer dfp.Close()
	if err != nil {
		return err
	}

	sfp, err := os.Open(r.Path)
	defer sfp.Close()
	if err != nil {
		return err
	}

	_, err = io.Copy(dfp,sfp)

	return err
}
