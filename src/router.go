package src

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/nathan-hello/htmx-template/src/components"
)

func Todo(w http.ResponseWriter, r *http.Request) {

	post := func() {
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		body := r.FormValue("body")

		if len(body) > 255 {
			w.WriteHeader(http.StatusBadRequest)
		}

	}

}

func Router() {
	http.HandleFunc("/static/css/tw-output.css", func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "text/css")
		http.ServeFile(res, req, "src/static/css/tw-output.css")
	})

	http.HandleFunc("/favicon.ico", func(res http.ResponseWriter, req *http.Request) {
		http.ServeFile(res, req, "src/static/favicon.ico")
	})

	http.HandleFunc("/white-bear.ico", func(res http.ResponseWriter, req *http.Request) {
		http.ServeFile(res, req, "src/static/white-bear.ico")
	})

	http.Handle("/", templ.Handler(components.Root()))

	// mime.AddExtensionType(".css", "text/css")
	http.ListenAndServe(":3000", nil)
}
