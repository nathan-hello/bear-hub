package routes

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/nathan-hello/htmx-template/src/components"
	"github.com/nathan-hello/htmx-template/src/db"
)

func UserProfile(w http.ResponseWriter, r *http.Request) {

	pathSegments := strings.Split(r.URL.Path, "/")
	if pathSegments[1] != "profile" {
		http.Redirect(w, r, "/profile", http.StatusSeeOther)

		return
	}

	requestedProfile := pathSegments[2]

	todos, err := db.Db().SelectTodosByUsername(r.Context(), requestedProfile)

	if err != nil {
		if err != sql.ErrNoRows {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	}

	p := components.ProfileProps{Username: requestedProfile, Todos: &todos}

	components.Profile(p).Render(r.Context(), w)
}
