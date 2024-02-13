package routes

import (
	"net/http"

	"bear-hub/examples/bear-hub/src/components"
)

func HandleRedirect(w http.ResponseWriter, r *http.Request, route string, status int, err error) {
	w.WriteHeader(status)
	http.Redirect(w, r, route, http.StatusSeeOther)
	if err != nil {
		components.AlertBox(err.Error()).Render(r.Context(), w)
	}

}
