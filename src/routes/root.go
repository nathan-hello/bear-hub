package routes

import (
	"context"
	"net/http"

	"github.com/a-h/templ"
	"github.com/nathan-hello/htmx-template/src/components"
)

func Alert(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/alert/signout" {
		w.Write([]byte("You've been signed out"))
	}
	if r.URL.Path == "/alert/unauthorized" {
		w.Write([]byte("You're not signed in"))
	}
}

func Root(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
	}
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	ctx := context.Background()

	response, err := templ.ToGoHTML(ctx, components.Root())

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Write([]byte(response))

}
