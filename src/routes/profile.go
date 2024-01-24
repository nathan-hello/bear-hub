package routes

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/nathan-hello/htmx-template/src/components"
	"github.com/nathan-hello/htmx-template/src/db"
	"github.com/nathan-hello/htmx-template/src/utils"
)

func convertInt64ToInt32(i []int64) []int32 {
	new := make([]int32, len(i))
	for _, v := range i {
		new = append(new, int32(v))
	}
	return new
}

func UserProfile(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()
	d, err := sql.Open("postgres", utils.Env().DB_URI)
	if err != nil {
		Redirect500(w, r)
		return
	}

	conn := db.New(d)

	pathSegments := strings.Split(r.URL.Path, "/")
	if pathSegments[1] != "profile" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	requestedProfile := pathSegments[2]

	todos, err := conn.SelectTodosByUsername(ctx, requestedProfile)

	if err != nil {
		if err != sql.ErrNoRows {
			fmt.Printf("err: %#v\n", err)
			Redirect500(w, r)
		}
	}

	components.Profile(&components.ProfileProps{
		Username: requestedProfile,
		Todos:    &todos,
	}).Render(r.Context(), w)
}
