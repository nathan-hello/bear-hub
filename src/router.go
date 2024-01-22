package src

import (
	"net/http"

	"github.com/nathan-hello/htmx-template/src/routes"
)

func staticGet(path string, filepath string) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path != path {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		http.ServeFile(w, r, filepath)
	}

}

func Router() {
	http.HandleFunc("/static/css/tw-output.css", func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "text/css")
		http.ServeFile(res, req, "src/static/css/tw-output.css")
	})

	http.HandleFunc("/favicon.ico", staticGet("/favicon.ico", "src/static/favicon.ico"))
	http.HandleFunc("/white-bear.ico", staticGet("/white-bear.ico", "src/static/white-bear.ico"))

	http.HandleFunc("/", routes.Root)
	http.HandleFunc("/todo", routes.Todo)
	http.HandleFunc("/signup", routes.SignUp)
	http.HandleFunc("/signin", routes.SignIn)
	// http.HandleFunc("/profile/", routes.UserProfile)
	http.HandleFunc("/404", routes.NotFound)
	http.HandleFunc("/500", routes.InternalServerError)
	http.HandleFunc("/testing/", routes.TestingEntryPoint)

	http.ListenAndServe(":3000", nil)
}
