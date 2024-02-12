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
		HandleRedirect(w, r, "/signin", http.StatusUnauthorized, utils.ErrBadLogin)
		return
	}
	conn, err := utils.Db()
	if err != nil {
		HandleRedirect(w, r, "/?500=/", http.StatusInternalServerError, utils.ErrDbConnection)
		return
	}

	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			HandleRedirect(w, r, "/?500=/", http.StatusInternalServerError, utils.ErrDbConnection)
			return
		}

		body := r.FormValue("body")

		if len(body) > 255 || len(body) < 3 {
			// this should be a send htmx div like in auth.go
			HandleRedirect(w, r, "/?500=/", http.StatusBadRequest, utils.ErrBadReqTodosBodyShort)
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
		err = conn.DeleteTodo(ctx, parsedId)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method == "GET" {
		todos, err := conn.SelectTodosByUsername(ctx, claims.Username)
		if err != nil {
			if err != sql.ErrNoRows {
				HandleRedirect(w, r, "/?500=/", http.StatusInternalServerError, utils.ErrDbSelectTodosByUser)
				return
			}
		}
		components.Todo(todos).Render(r.Context(), w)
	}

}
