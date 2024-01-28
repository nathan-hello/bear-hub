package src

import (
	"context"
	"log"
	"net/http"
	"slices"
	"time"

	"github.com/justinas/alice"
	"github.com/nathan-hello/htmx-template/src/routes"
	"github.com/nathan-hello/htmx-template/src/utils"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			log.Printf("IP: %v, ROUTE REQUESTED: %v, RESPONSE TIME: %v\n", r.URL.User, r.URL.Path, time.Since(start))
		}()
		next.ServeHTTP(w, r)
	})
}
func AllowMethods(methods ...string) alice.Constructor {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !slices.Contains(methods, r.Method) {
				w.WriteHeader(http.StatusMethodNotAllowed)
				w.Write([]byte("Method not allowed"))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// This returned an alice.Constructor instead of
// being an alice.Constructor because it requires an argument (path string).
func CreateHeader(key string, value string) alice.Constructor {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add(key, value)
			next.ServeHTTP(w, r)
		})
	}
}

// This returned an alice.Constructor instead of
// being an alice.Constructor because it requires an argument (path string).
func RejectSubroute(path string) alice.Constructor {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != path {
				routes.Redirect404(w, r)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func InjectClaimsOnValidToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		access, ok := utils.ValidateJwtOrDelete(w, r)
		if !ok {
			next.ServeHTTP(w, r)
			return
		}

		claims, err := utils.ParseToken(access)
		if err != nil {
			log.Print("ERR: parsetoken:", access, err)
			next.ServeHTTP(w, r)
			return
		}

		// if claims is nil, it doesn't matter because
		// we have to do a type assertion whenever we use it anyways
		// and that will check if the type is ok
		var claimsObj *utils.CustomClaims = claims
		newCtx := context.WithValue(r.Context(), utils.ClaimsContextKey, claimsObj)

		next.ServeHTTP(w, r.WithContext(newCtx))
	})
}
