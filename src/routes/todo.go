package routes

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/nathan-hello/htmx-template/src/components"
	"github.com/nathan-hello/htmx-template/src/db"
	"github.com/nathan-hello/htmx-template/src/utils"
)

func Todo(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" && r.Method != "DELETE" && r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// var response template.HTML
	ctx := context.Background()
	d, err := sql.Open("postgres", utils.Env().DB_URI)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	conn := db.New(d)

	if r.Method == "POST" {

		access, ok := utils.ValidateJwtOrDelete(w, r)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		claims, err := utils.ParseToken(access)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		body := r.FormValue("body")

		if len(body) > 255 || len(body) < 3 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		row, err := conn.InsertTodo(ctx, db.InsertTodoParams{Body: body, Username: claims.Username})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		components.TodoRow(&row).Render(r.Context(), w)
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
		access, ok := utils.ValidateJwtOrDelete(w, r)
		if !ok {
			RedirectToSignIn(w, r)
			return
		}

		claims, err := utils.ParseToken(access)
		if err != nil {
			RedirectToSignIn(w, r)
			return
		}
		todos, err := conn.SelectTodosByUsername(ctx, claims.Username)
		components.Todo(&todos).Render(r.Context(), w)
	}

}
