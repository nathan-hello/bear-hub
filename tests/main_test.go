package test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/lib/pq"
	"github.com/nathan-hello/htmx-template/src/config"
	"github.com/nathan-hello/htmx-template/src/sqlc"
)

func TestDatabaseConnection(t *testing.T) {
	ctx := context.Background()

	db, err := sql.Open("postgres", config.Get().DB_URI)

	if err != nil {
		panic(err)
	}

	f := sqlc.New(db)

	rows, err := f.AllTodos(ctx, 10)

	if err != nil {
		panic(err)
	}

	fmt.Printf("rows of todos: %#v", rows)

}
