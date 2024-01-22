// package routes
//
// import (
//
//	"context"
//	"database/sql"
//	"net/http"
//	"strings"
//
//	"github.com/nathan-hello/htmx-template/src/components"
//	"github.com/nathan-hello/htmx-template/src/db"
//	"github.com/nathan-hello/htmx-template/src/utils"
//
// )
//
//	func convertInt64ToInt32(i []int64) []int32 {
//		new := make([]int32, len(i))
//		for _, v := range i {
//			new = append(new, int32(v))
//		}
//		return new
//	}
//
// func UserProfile(w http.ResponseWriter, r *http.Request) {
//
//		ctx := context.Background()
//		d, err := sql.Open("postgres", utils.Env().DB_URI)
//		if err != nil {
//			redirectServerError(w, r)
//			return
//		}
//
//		conn := db.New(d)
//
//		pathSegments := strings.Split(r.URL.Path, "/")
//		if pathSegments[1] != "profile" {
//			w.WriteHeader(http.StatusNotFound)
//			return
//		}
//
//		requestedProfile := pathSegments[2]
//
//		row, err := conn.SelectProfileByUsername(ctx, requestedProfile)
//
//		if err != nil {
//			redirectNotFound(w, r)
//			return
//		}
//
//		 todos, err := conn.SelectTodosByIds(
//		 	ctx,
//		 	sequel.SelectTodosByIdsParams{
//		 		Column1: convertInt64ToInt32(row.Todos),
//		 		Limit:   10,
//		 	})
//
//		 if err != nil {
//		 	redirectServerError(w, r)
//		 	return
//		 }
//
//		components.Profile(&components.ProfileProps{
//			Username: row.Username,
//			Todos:    &[]db.Todo{},
//		}).Render(r.Context(), w)
//	}
package routes

func Asdf() int {
	return 3
}
