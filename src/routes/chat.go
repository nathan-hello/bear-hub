package routes

import (
	"log"
	"net/http"

	gws "github.com/gorilla/websocket"
	"github.com/nathan-hello/htmx-template/src/components"
)

var upgrader = gws.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ChatSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	log.Println("HIT!")
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	for {
		_, msg, err := conn.ReadMessage()
		log.Printf("1: msg: %#v, err: %#v\n", string(msg), err)
		msg = []byte("123123123123")
		if err != nil {
			log.Println(err)
			return
		}
		if err := conn.WriteMessage(gws.TextMessage, msg); err != nil {
			log.Println(err)
			return
		}
		log.Printf("2: msg: %s, err: %#v\n", string(msg), err)
	}

}

func Chat(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		components.ChatRoomRoot().Render(r.Context(), w)
		return
	}
}

// func postChat(w http.ResponseWriter, r *http.Request) {
// 	claims, ok := r.Context().Value(utils.ClaimsContextKey).(utils.CustomClaims)
// 	if !ok {
// 		w.WriteHeader(http.StatusUnauthorized)
// 		return
// 	}
//
// 	if err := r.ParseForm(); err != nil {
// 		w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}
//
// 	conn, err := utils.Db()
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}
//
// 	t, err := conn.InsertChatroom(
// 		r.Context(),
// 		db.InsertChatroomParams{
// 			Name:    r.FormValue("name"),
// 			Creator: claims.Username,
// 		})
//
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 	}
//
// }
