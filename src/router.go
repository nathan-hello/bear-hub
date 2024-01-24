package src

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/nathan-hello/htmx-template/src/routes"
	"github.com/nathan-hello/htmx-template/src/utils"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

// Logging logs all requests with its path and the time it took to process
func Logging() Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Do middleware things
			start := time.Now()
			defer func() { log.Println(r.URL.Path, time.Since(start)) }()
			// Call the next middleware/handler in chain
			f(w, r)
		}
	}
}

// Method ensures that url can only be requested with a specific method, else returns a 400 Bad Request
func Method(m string) Middleware {
	// Create a new Middleware
	return func(f http.HandlerFunc) http.HandlerFunc {
		// Define the http.HandlerFunc
		return func(w http.ResponseWriter, r *http.Request) {
			// Do middleware things
			if r.Method != m {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			// Call the next middleware/handler in chain
			f(w, r)
		}
	}
}

func ProtectedRoute() Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {

			access, ok := utils.ValidateJwtOrDelete(w, r)

			if !ok {
				if r.Method == "GET" {
					routes.RedirectToSignIn(w, r)
				} else {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
			}

			claims, err := utils.ParseToken(access)
			if err != nil {
				if r.Method == "GET" {
					routes.RedirectToSignIn(w, r)
				} else {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
			}

			ctx := context.WithValue(r.Context(), "claims", claims)
			r.WithContext(ctx)

			f(w, r)
		}
	}
}

func Chain(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		f = m(f)
	}
	return f
}

func staticGet(path string, filepath string) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path != path {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		http.ServeFile(w, r, filepath)
	}

}

func Router() {
	http.HandleFunc("/static/css/tw-output.css", func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "text/css")
		http.ServeFile(res, req, "src/static/css/tw-output.css")
	})

	http.HandleFunc("/favicon.ico", staticGet("/favicon.ico", "src/static/favicon.ico"))
	http.HandleFunc("/white-bear.ico", staticGet("/white-bear.ico", "src/static/white-bear.ico"))

	http.HandleFunc("/", routes.Root)
	http.HandleFunc("/todo", Chain(routes.Todo, Logging(), ProtectedRoute()))
	http.HandleFunc("/signup", routes.SignUp)
	http.HandleFunc("/signin", routes.SignIn)
	http.HandleFunc("/alert/", routes.Alert)
	http.HandleFunc("/profile/", routes.UserProfile)
	http.HandleFunc("/c/", routes.MicroComponents)

	http.ListenAndServe(":3000", nil)
}
