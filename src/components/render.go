package components

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"
	"github.com/nathan-hello/htmx-template/src/db"
	"github.com/nathan-hello/htmx-template/src/utils"
)

func getTodos() ([]db.Todo, error) {
	ctx := context.Background()
	d, err := sql.Open("postgres", utils.Env().DB_URI)

	if err != nil {
		return nil, err
	}

	todosTable := db.New(d)
	rows, err := todosTable.SelectTodos(ctx, 99)

	if err != nil {
		return nil, err
	}

	return rows, err
}

func formatTime(t *time.Time) string {
	return t.Format(time.RFC822)
}
