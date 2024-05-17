package routes

import (
	"net/http"

	"github.com/nathan-hello/htmx-template/src/auth"
	"github.com/nathan-hello/htmx-template/src/components"
)

func MicroComponents(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	if r.URL.Path == "/c/alert/signout" {
		components.AlertBox("You have been signed out").Render(r.Context(), w)
	}

	if r.URL.Path == "/c/alert/unauthorized" {
		components.AlertBox("You're not logged in").Render(r.Context(), w)
	}

	if r.URL.Path == "/c/alert/404" {
		components.AlertBox("404 Not Found").Render(r.Context(), w)
	}

	if r.URL.Path == "/c/alert/500" {
		components.AlertBox("500 Internal Server Error").Render(r.Context(), w)
	}
}

func Root(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value(auth.ClaimsContextKey).(*auth.CustomClaims)
        state := components.ClientState{
                IsAuthed: ok,
        }
	components.Root(state).Render(r.Context(), w)
}
