package routes

import (
	"html/template"
	"net/http"
)

func Root(w http.ResponseWriter, _ *http.Request) {
	tmpl, err := template.ParseFiles("src/static/templates/root.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Router() {
	http.HandleFunc("/static/css/output.css", func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "text/css")
		http.ServeFile(res, req, "src/static/css/output.css")
	})

	http.HandleFunc("/", Root)

	http.ListenAndServe(":3000", nil)
}
