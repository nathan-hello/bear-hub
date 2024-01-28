package routes

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/nathan-hello/htmx-template/src/components"
	"github.com/nathan-hello/htmx-template/src/utils"
)

func MicroComponents(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	if r.URL.Path == "/c/alert/signout" {
		components.NotificationBox("You have been signed out").Render(r.Context(), w)
	}

	if r.URL.Path == "/c/alert/unauthorized" {
		components.NotificationBox("You're not logged in").Render(r.Context(), w)
	}

	if r.URL.Path == "/c/alert/404" {
		components.NotificationBox("404 Not Found").Render(r.Context(), w)
	}

	if r.URL.Path == "/c/alert/500" {
		components.NotificationBox("500 Internal Server Error").Render(r.Context(), w)
	}

	if r.URL.Path == "/c/isloggedin" {
		claims, ok := r.Context().Value(utils.ClaimsContextKey).(utils.CustomClaims)
		if ok {
			p := fmt.Sprintf("/profile/%v", claims.Username)
			components.NavbarLink(templ.SafeURL(p), claims.Username).Render(r.Context(), w)
			return
		} else {
			components.NavbarLink("/signin", "Sign In").Render(r.Context(), w)
		}
	}
}

func Root(w http.ResponseWriter, r *http.Request) {
	components.Root().Render(r.Context(), w)
}
