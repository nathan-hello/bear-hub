package src

import (
	"context"
	"log"
	"net/http"
	"slices"
	"time"

	"github.com/justinas/alice"
	"github.com/nathan-hello/htmx-template/src/auth"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			log.Printf("IP: %v, ROUTE REQUESTED: %v, RESPONSE TIME: %v\n", r.RemoteAddr, r.URL.Path, time.Since(start))
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
// being an alice.Constructor because it requires an argument (key, value string).
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
				http.NotFound(w, r)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func InjectClaimsOnValidToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		access, ok := auth.ValidateJwtOrDelete(w, r)
		if !ok {
			next.ServeHTTP(w, r)
			return
		}

		claims, err := auth.ParseToken(access)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

                if claims == nil {
                        log.Println("claims was nil")
			next.ServeHTTP(w, r)
			return

                }
                wrapReq := r.WithContext(context.WithValue(r.Context(), auth.ClaimsContextKey, claims))
		next.ServeHTTP(w, wrapReq)
	})
}
