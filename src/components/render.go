package components

import (
	"time"

	"github.com/nathan-hello/htmx-template/src/db"
)
type LayoutParams struct {
	TabTitle string
	NavTitle string
}

type ProfileProps struct {
	Username string
	Todos    *[]db.Todo
}

func formatTime(t *time.Time) string {
	return t.Format(time.RFC822)
}

type ClientState struct {
	IsAuthed bool
}
