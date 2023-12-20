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

func Todos(w http.ResponseWriter, _ *http.Request) {
	tmpl, err := template.ParseFiles("src/static/templates/todos.html")

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
	http.HandleFunc("/static/css/tw-output.css", func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "text/css")
		http.ServeFile(res, req, "src/static/css/tw-output.css")
	})

	http.HandleFunc("/favicon.ico", func(res http.ResponseWriter, req *http.Request) {
		http.ServeFile(res, req, "src/static/favicon.ico")
	})

	http.HandleFunc("/", Root)

	// mime.AddExtensionType(".css", "text/css")
	http.ListenAndServe(":3000", nil)
}
