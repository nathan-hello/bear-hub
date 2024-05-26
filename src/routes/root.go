package routes

import (
	"net/http"

	"github.com/nathan-hello/htmx-template/src/auth"
	"github.com/nathan-hello/htmx-template/src/components"
)

func Root(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value(auth.ClaimsContextKey).(*auth.CustomClaims)
        state := components.ClientState{
                IsAuthed: ok,
        }
	components.Root(state).Render(r.Context(), w)
}
