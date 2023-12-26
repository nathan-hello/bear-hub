package routes

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

	if r.URL.Path != "/todo" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if r.Method != "POST" && r.Method != "DELETE" && r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var response template.HTML
	ctx := context.Background()
	db, err := sql.Open("postgres", utils.Env().DB_URI)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	conn := sqlc.New(db)

	if r.Method == "POST" {

		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		body := r.FormValue("body")

		if len(body) > 255 || len(body) < 3 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		row, err := conn.InsertTodo(ctx, body)
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

	if r.Method == "DELETE" {
		id := r.URL.Query().Get("id")
		parsedId, err := strconv.ParseInt(id, 10, 0)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		conn.DeleteTodo(ctx, parsedId)
		w.WriteHeader(http.StatusOK)

	}

	if r.Method == "GET" {

		response, err := templ.ToGoHTML(ctx, components.Todo())

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write([]byte(response))
		return

	}

}
