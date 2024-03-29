package routes

import (
	"net/http"

	"github.com/nathan-hello/htmx-template/src/components"
)

func HandleRedirect(w http.ResponseWriter, r *http.Request, route string, err error) {
	http.Redirect(w, r, route, http.StatusSeeOther)
	if err != nil {
		components.AlertBox(err.Error()).Render(r.Context(), w)
	}

}
