package components

import (
	"time"

	"bear-hub/examples/bear-hub/examples/bear-hub/src/db"
	"bear-hub/examples/bear-hub/examples/bear-hub/src/utils"
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
	if s == nil {
		return nil
	}
	e := map[string]FieldError{}
	for _, v := range *s {
		e[v.Field] = FieldError{
			BorderColor: "bg-red-500",
			Value:       v.Value,
			Err:         v.Err.Error(),
		}
	}
	s = nil
	return e
}
