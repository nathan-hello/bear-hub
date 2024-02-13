package routes

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	gws "github.com/gorilla/websocket"
	"github.com/nathan-hello/htmx-template/examples/bear-hub/examples/bear-hub/src/components"
)

var upgrader = gws.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Manager struct {
	clients map[*gws.Conn]bool
	lock    sync.Mutex
}

func (m *Manager) AddClient(c *gws.Conn) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.clients[c] = true
}

func (m *Manager) RemoveClient(c *gws.Conn) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if _, ok := m.clients[c]; ok {
		delete(m.clients, c)
		c.Close()
	}
}

func (m *Manager) BroadcastMessage(message []byte) {
	m.lock.Lock()
	defer m.lock.Unlock()
	for c, _ := range m.clients {
		if err := c.WriteMessage(gws.TextMessage, message); err != nil {
			log.Println(err)
			delete(m.clients, c)
			c.Close()
		}
	}
}

var manager = Manager{
	clients: make(map[*gws.Conn]bool),
}

type ChatParse struct {
	Text string `json:"msg-input"`
}

type ChatResponse struct {
	User  string
	Color string
	Text  string
}

func ChatSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	manager.AddClient(conn)
	defer manager.RemoveClient(conn)

	for {
		_, clientMsg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("clientMsg: %s, err: %#v\n", clientMsg, err)

		t := &ChatParse{}
		json.Unmarshal(clientMsg, &t)
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("unmarshal: %#v, err: %#v\n", t, err)

		resp := ChatResponse{
			User:  "rewq",
			Color: "bg-blue-200",
			Text:  t.Text,
		}

		var buffMsg bytes.Buffer
		components.ChatMessage(resp.Text).Render(r.Context(), &buffMsg)
		log.Printf("buffMsg: %s, err: %#v\n", buffMsg.Bytes(), err)

		manager.BroadcastMessage(buffMsg.Bytes())

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
