package test

import (
	"fmt"

	_ "github.com/lib/pq"
	"github.com/nathan-hello/htmx-template/src/config"
	todos "github.com/nathan-hello/htmx-template/src/db"
)

func Database() {

	db, ctx, err := todos.Connection(config.Get().DB_URI)

	if err != nil {
		panic(err)
	}

	rows, err := db.AllTodos(*ctx, 99)

	if err != nil {
		panic(err)
	}

	fmt.Printf("rows of todos: %#v", rows)

}
