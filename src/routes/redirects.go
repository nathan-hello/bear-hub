package routes

import (
	"fmt"
	"net/http"

	"github.com/nathan-hello/htmx-template/src/components"
)

func redirectNotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("HX-Redirect", "404")
	fmt.Printf("%#v\n", r)
}

func redirectServerError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("HX-Redirect", "500")
	fmt.Printf("%#v\n", r)
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/404" {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("HX-Redirect", "404")
		return
	}

	if r.Method == "GET" {
		components.Error("404 Not Found").Render(r.Context(), w)
	}
}

func InternalServerError(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/500" {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("HX-Redirect", "404")
		return
	}

	if r.Method == "GET" {
		components.Error("500 Internal Server Error").Render(r.Context(), w)
	}
}
