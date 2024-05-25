package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nathan-hello/htmx-template/src/utils"
)

var db *Queries

func DbInit() error {
	var d, err = sql.Open("sqlite3", utils.Env().DB_URI)
	if err != nil {
		return err
	}
	db = New(d)
	return nil
}

func Db() *Queries {
	return db
}
