package routes

import (
	"log"
	"net/http"

	"github.com/nathan-hello/htmx-template/src/auth"
	"github.com/nathan-hello/htmx-template/src/components"
)

func Root(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.ClaimsContextKey).(*auth.CustomClaims)
        state := components.ClientState{
                IsAuthed: ok,
        }
        log.Println(claims)
	components.Root(state).Render(r.Context(), w)
}
