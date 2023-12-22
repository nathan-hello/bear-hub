package src

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/a-h/templ"
	_ "github.com/lib/pq"
	"github.com/nathan-hello/htmx-template/src/components"
	"github.com/nathan-hello/htmx-template/src/sqlc"
)

func Router() {
	http.HandleFunc("/static/css/tw-output.css", func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "text/css")
		http.ServeFile(res, req, "src/static/css/tw-output.css")
	})

	http.HandleFunc("/favicon.ico", func(res http.ResponseWriter, req *http.Request) {
		http.ServeFile(res, req, "src/static/favicon.ico")
	})

	ctx := context.Background()
	db, err := sql.Open("postgres", Env().DB_URI)
	if err != nil {
		panic(err)
	}

	todosTable := sqlc.New(db)
	rows, err := todosTable.AllTodos(ctx, 99)

	if err != nil {
		panic(err)
	}

	http.Handle("/", templ.Handler(components.Root(rows)))

	// mime.AddExtensionType(".css", "text/css")
	http.ListenAndServe(":3000", nil)
}
