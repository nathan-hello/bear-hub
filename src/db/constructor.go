package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	_ "github.com/lib/pq"
	"github.com/nathan-hello/htmx-template/src/db/output/pg"
	"github.com/nathan-hello/htmx-template/src/db/output/sqlite"
	"github.com/nathan-hello/htmx-template/src/utils"
)

// The rest of the application uses this file to interact with the database.
// Two reasons two abstract in this way:
//      1: Testing
//      2: Hiding implementation for type coersion.
// In the event that I do switch from postgres to sqlite3, having this
// be the interface for the database means I don't have to run around the entire
// project to change the types. This has bitten me from going from postgres to
// sqlite3, so I figure if I'm going to put in the effort to change much of the
// types necessary, I might as well make an interface.
//
// This is the constructor for the wrapping of sqlc calls. If I could make an
// interface that encompasses everything I want, I could do that. Not sure if it's
// worth. At the least, we're wrapping around multiple sqlc Querier interfaces
// Though we use *sqlc.Queries instead of sqlc.Querier because that's a more accurate
// type to what we're actually working with I think.

var pgDb pg.Querier
var slDb sqlite.Querier

type dbDelete struct {
	pgDb *pg.Queries
	slDb *sqlite.Queries
}
type dbUpdate struct {
	pgDb *pg.Queries
	slDb *sqlite.Queries
}
type dbInsert struct {
	pgDb *pg.Queries
	slDb *sqlite.Queries
}
type dbSelect struct {
	pgDb *pg.Queries
	slDb *sqlite.Queries
}

var dbType = utils.Env().DB_TYPE
var Delete dbDelete
var Update dbUpdate
var Insert dbInsert
var Select dbSelect

func pgInit() (*pg.Queries, error) {
	var d, err = pgx.Connect(context.Background(), utils.Env().PG_DB_URI)
	if err != nil {
		return nil, err
	}
	return pg.New(d), nil
}

func sqliteInit() (*sqlite.Queries, error) {
	d, err := sql.Open("sqlite3", utils.Env().SQLITE_DB_URI)
	if err != nil {
		return nil, err
	}
	q := sqlite.New(d)
	return q, nil
}
func InitDb() error {
	dbType := utils.Env().DB_TYPE
	switch dbType {
	case "postgres":
		d, err := pgInit()
		pgDb = d
		Delete = dbDelete{pgDb: d}
		Update = dbUpdate{pgDb: d}
		Insert = dbInsert{pgDb: d}
		Select = dbSelect{pgDb: d}
		return err

	case "sqlite3":
		d, err := sqliteInit()
		slDb = d
		Delete = dbDelete{slDb: d}
		Update = dbUpdate{slDb: d}
		Insert = dbInsert{slDb: d}
		Select = dbSelect{slDb: d}
		return err
	}

	return errors.New("DB_TYPE not found")
}

func DbSanityCheck() error {
	if pgDb == nil && slDb == nil {
		return errors.New("both pgDb and slDb are nil")
	}
	if Delete.pgDb == nil && Delete.slDb == nil {
		return fmt.Errorf("delete is uninitiated: %#v", Delete)
	}
	if Update.pgDb == nil && Update.slDb == nil {
		return fmt.Errorf("update is uninitiated: %#v", Update)
	}
	if Insert.pgDb == nil && Insert.slDb == nil {
		return fmt.Errorf("insert is uninitiated: %#v", Insert)
	}
	if Select.pgDb == nil && Select.slDb == nil {
		return fmt.Errorf("select is uninitiated: %#v", Select)
	}
	return nil
}

type DbError struct {
	slErr error
	pgErr error
}

func (d *DbError) Error() string {
	return fmt.Sprintf("database error: %#v", d)
}

func (d *DbError) Unwrap() error {
	if d.slErr == nil && d.pgErr == nil {
		return nil
	}
	return d
}
