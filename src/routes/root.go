package routes

import (
	"net/http"

	"github.com/nathan-hello/htmx-template/src/components"
	"github.com/nathan-hello/htmx-template/src/utils"
)

func Root(w http.ResponseWriter, r *http.Request) {
        state := utils.GetClientState(r)
	components.Root(state).Render(r.Context(), w)
}
