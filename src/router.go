package src

import (
	"context"
	"database/sql"
	"html/template"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/nathan-hello/htmx-template/src/components"
	"github.com/nathan-hello/htmx-template/src/sqlc"
	"github.com/nathan-hello/htmx-template/src/utils"
)

func Todo(w http.ResponseWriter, r *http.Request) {

	var response template.HTML
	ctx := context.Background()

	if r.URL.Path != "/todo" || r.Method != "POST" && r.Method != "DELETE" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	db, err := sql.Open("postgres", utils.Env().DB_URI)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	todosTable := sqlc.New(db)

	post := func() {

		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		body := r.FormValue("body")

		if len(body) > 255 || len(body) < 3 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		row, err := todosTable.InsertTodo(ctx, body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		response, err = templ.ToGoHTML(ctx, components.TodoRow(&row))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write([]byte(response))
		return
	}

	del := func() {
		id := r.URL.Query().Get("id")
		parsedId, err := strconv.ParseInt(id, 10, 0)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		todosTable.DeleteTodo(ctx, parsedId)
		w.WriteHeader(http.StatusOK)

	}

	if r.Method == "POST" {
		post()
		return
	}

	if r.Method == "DELETE" {
		del()
		return
	}

}

func Root(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" || r.Method != "GET" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	ctx := context.Background()

	response, err := templ.ToGoHTML(ctx, components.Root())

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Write([]byte(response))
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

	http.HandleFunc("/", Root)
	http.HandleFunc("/todo", Todo)

	// mime.AddExtensionType(".css", "text/css")
	http.ListenAndServe(":3000", nil)
}
