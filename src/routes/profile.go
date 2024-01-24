package routes

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/nathan-hello/htmx-template/src/components"
	"github.com/nathan-hello/htmx-template/src/utils"
)

func UserProfile(w http.ResponseWriter, r *http.Request) {
	conn, err := utils.Db()
	if err != nil {
		Redirect500(w, r)
		return
	}

	pathSegments := strings.Split(r.URL.Path, "/")
	if pathSegments[1] != "profile" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	requestedProfile := pathSegments[2]

	todos, err := conn.SelectTodosByUsername(r.Context(), requestedProfile)

	if err != nil {
		if err != sql.ErrNoRows {
			fmt.Printf("err: %#v\n", err)
			Redirect500(w, r)
			return
		}
	}

	p := components.ProfileProps{Username: requestedProfile, Todos: &todos}

	components.Profile(&p).Render(r.Context(), w)
}
