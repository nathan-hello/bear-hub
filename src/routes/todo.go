package routes

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/nathan-hello/htmx-template/src/components"
	"github.com/nathan-hello/htmx-template/src/db"
	"github.com/nathan-hello/htmx-template/src/utils"
)

func Todo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := r.Context().Value(utils.ClaimsContextKey).(*utils.CustomClaims)
	if !ok {
		Redirect401(w, r)
		return
	}
	conn, err := utils.Db()
	if err != nil {
		Redirect500(w, r)
		return
	}

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
		return
	}

	if r.Method == "GET" {
		claims, ok := r.Context().Value(utils.ClaimsContextKey).(utils.CustomClaims)
		if !ok {
			Redirect500(w, r)
			return
		}

		todos, err := conn.SelectTodosByUsername(ctx, claims.Username)
		if err != nil {
			if err != sql.ErrNoRows {
				Redirect500(w, r)
				return
			}
		}
		components.Todo(todos).Render(r.Context(), w)
	}

}
