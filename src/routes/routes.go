package routes

import (
	"html/template"
	"net/http"
)

type resource struct {
	url  string
	real string
}

var Resources = map[string]resource{
	"output.css": {
		url:  "/static/css/output.css",
		real: "src/static/css/output.css",
	},
	"root.html": {
		url:  "/",
		real: "src/static/templates/root.html",
	},
}

func Root(w http.ResponseWriter, _ *http.Request) {
	tmpl, err := template.ParseFiles(Resources["root.html"].real)

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
	http.HandleFunc(Resources["output.css"].url, func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "text/css")
		http.ServeFile(res, req, Resources["output.css"].real)
	})

	http.HandleFunc("/", Root)

	http.ListenAndServe(":3000", nil)
}
