package components

import "github.com/nathan-hello/htmx-template/src/db"

type ProfileProps struct {
	Username string
	Todos    *[]db.Todo
}
