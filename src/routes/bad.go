package routes

import (
	"net/http"

	"github.com/nathan-hello/htmx-template/src/components"
)

func NotFound(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/404" {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("HX-Redirect", "404")
		return
	}

	if r.Method == "GET" {
		components.NotFound().Render(r.Context(), w)
	}
}
