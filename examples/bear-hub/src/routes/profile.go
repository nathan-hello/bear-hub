package routes

import (
	"database/sql"
	"net/http"
	"strings"

	"bear-hub/examples/bear-hub/src/components"
	"bear-hub/examples/bear-hub/src/utils"
)

func UserProfile(w http.ResponseWriter, r *http.Request) {
	conn, err := utils.Db()
	if err != nil {
		HandleRedirect(w, r, "/", http.StatusInternalServerError, utils.ErrDbConnection)
		return
	}

	pathSegments := strings.Split(r.URL.Path, "/")
	if pathSegments[1] != "profile" {
		HandleRedirect(w, r, "/profile", http.StatusNotFound, utils.ErrProfileNotFound)
		return
	}

	requestedProfile := pathSegments[2]

	todos, err := conn.SelectTodosByUsername(r.Context(), requestedProfile)

	if err != nil {
		if err != sql.ErrNoRows {
			HandleRedirect(w, r, "/", http.StatusInternalServerError, utils.ErrDbConnection)
		}
	}

	p := components.ProfileProps{Username: requestedProfile, Todos: &todos}

	components.Profile(p).Render(r.Context(), w)
}
