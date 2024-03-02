package routes

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/nathan-hello/htmx-template/src/components"
	"github.com/nathan-hello/htmx-template/src/utils"
)

func UserProfile(w http.ResponseWriter, r *http.Request) {
	conn, err := utils.Db()
	if err != nil {
		HandleRedirect(w, r, "/", utils.ErrDbConnection)
		return
	}

	pathSegments := strings.Split(r.URL.Path, "/")
	if pathSegments[1] != "profile" {
		HandleRedirect(w, r, "/profile", utils.ErrProfileNotFound)
		return
	}

	requestedProfile := pathSegments[2]

	todos, err := conn.SelectTodosByUsername(r.Context(), requestedProfile)

	if err != nil {
		if err != sql.ErrNoRows {
			HandleRedirect(w, r, "/", utils.ErrDbConnection)
		}
	}

	p := components.ProfileProps{Username: requestedProfile, Todos: &todos}

	components.Profile(p).Render(r.Context(), w)
}
