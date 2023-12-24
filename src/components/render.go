package components

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/nathan-hello/htmx-template/src/utils"
	"github.com/nathan-hello/htmx-template/src/sqlc"
)

func getTodos() ([]sqlc.Todo, error) {
	ctx := context.Background()
	db, err := sql.Open("postgres", utils.Env().DB_URI)

	if err != nil {
		return nil, err
	}

	todosTable := sqlc.New(db)
	rows, err := todosTable.AllTodos(ctx, 99)

	if err != nil {
		return nil, err
	}

	return rows, err
}
