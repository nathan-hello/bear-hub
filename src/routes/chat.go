package routes

import (
	"net/http"

	"github.com/nathan-hello/htmx-template/src/components"
	"github.com/nathan-hello/htmx-template/src/db"
	"github.com/nathan-hello/htmx-template/src/utils"
)

func ChatSubRouter(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		components.ChatRoomRoot().Render(r.Context(), w)
		return
	}

	if r.Method == "POST" {
		postChat(w, r)
		return
	}

	if r.Method == "DELETE" {

	}

}

func postChat(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(utils.ClaimsContextKey).(utils.CustomClaims)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	conn, err := utils.Db()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	t, err := conn.InsertChatroom(
		r.Context(),
		db.InsertChatroomParams{
			Name:    r.FormValue("name"),
			Creator: claims.Username,
		})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	RedirectToChatroom(w, r, t)
}

func deleteChat(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(utils.ClaimsContextKey).(utils.CustomClaims)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	conn, err := utils.Db()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	

}
