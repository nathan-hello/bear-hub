package components

import (
	"time"

	"github.com/nathan-hello/htmx-template/src/db"
	"github.com/nathan-hello/htmx-template/src/utils"
)

type ProfileProps struct {
	Username string
	Todos    *[]db.Todo
}

func formatTime(t *time.Time) string {
	return t.Format(time.RFC822)
}

type FieldError struct {
	BorderColor string
	Value       string
	Err         string
}

func RenderAuthError(s *[]utils.AuthError) map[string]FieldError {
	asdf := map[string]FieldError{}

	for _, v := range utils.AllFields {
		asdf[v] = FieldError{BorderColor: "bg-blue-500", Value: ""}
	}

	for _, v := range *s {
		asdf[v.Field] = FieldError{
			BorderColor: "bg-red-500",
			Value:       v.Value,
			Err:         v.Err.Error(), //
		}
	}

	return asdf
}
