package routes

import (
	"context"
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

	if r.URL.Path == "/c/isloggedin" {
		if token, ok := utils.ValidateJwtOrDelete(w, r); ok {
			claims, err := utils.ParseToken(token)
			if err != nil {
				utils.DeleteJwtCookies(w)
				Redirect500(w, r)
				return
			}
			p := fmt.Sprintf("/profile/%v", claims.Username)
			components.NavbarLink(templ.SafeURL(p), claims.Username).Render(r.Context(), w)
			return
		} else {
			components.NavbarLink("/signin", "Sign In").Render(r.Context(), w)
		}
	}
}

func Root(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
	}
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	ctx := context.Background()

	response, err := templ.ToGoHTML(ctx, components.Root())

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Write([]byte(response))

}
