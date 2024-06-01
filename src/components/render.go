package components

import (
	"time"
	"unicode"

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

func formatTime(t time.Time) string {
	return t.Format(time.RFC822)
}

func sentenceizeString(s string) string {
	if s == "" {
		return ""
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	runes = append(runes, '.')
	return string(runes)
}
