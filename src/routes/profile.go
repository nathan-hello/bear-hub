package routes

import (
	"context"
	"database/sql"
	"net/http"
	"strings"

	"github.com/a-h/templ"
	"github.com/nathan-hello/htmx-template/src/components"
	"github.com/nathan-hello/htmx-template/src/db"
	"github.com/nathan-hello/htmx-template/src/utils"
)

func UserProfile(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()
	d, err := sql.Open("postgres", utils.Env().DB_URI)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	conn := db.New(d)

	pathSegments := strings.Split(r.URL.Path, "/")
	if pathSegments[1] != "profile" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	requestedProfile := pathSegments[2]

	row, err := conn.SelectProfileByUsername(ctx, requestedProfile)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	response, err := templ.ToGoHTML(ctx, components.Profile(&row))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write([]byte(response))
	return

}
