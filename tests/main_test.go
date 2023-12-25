package test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/lib/pq"
	"github.com/nathan-hello/htmx-template/src/sqlc"
	"github.com/nathan-hello/htmx-template/src/utils"
)

func TestDatabaseConnection(t *testing.T) {
	ctx := context.Background()

	db, err := sql.Open("postgres", utils.Env().DB_URI)

	if err != nil {
		panic(err)
	}

	f := sqlc.New(db)

	rows, err := f.SelectTodos(ctx, 10)

	if err != nil {
		panic(err)
	}

	fmt.Printf("rows of todos: %#v\n", rows)

	row, err := f.InsertTodo(ctx, "gotest 1")

	if err != nil {
		panic(err)
	}

	fmt.Printf("inserted row: %#v\n", row)

	// f.DeleteTodo(ctx, row.ID)

	// fmt.Printf("deleted row: %#v\n", row)

}
