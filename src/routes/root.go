package routes

import (
	"context"
	"net/http"

	"github.com/a-h/templ"
	"github.com/nathan-hello/htmx-template/src/components"
)

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
