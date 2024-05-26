package components

import (
	"slices"
	"time"

	"github.com/nathan-hello/htmx-template/src/db"
)

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

type AuthState struct {
	ClientState
	Username    string
	UsernameErr string
	Email       string
	EmailErr    string
	PassErr     string
	PassConfErr string
}

func (a *AuthState) RenderErrs() []string {
	errs := []string{a.UsernameErr, a.EmailErr, a.PassErr, a.PassConfErr}
	for i, v := range errs {
		if v == "" {
			errs = slices.Delete(errs, i, i+1)
		}
	}
	return errs
}

type LayoutParams struct {
	TabTitle string
	NavTitle string
}
