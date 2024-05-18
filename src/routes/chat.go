package routes

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	gws "github.com/gorilla/websocket"
	"github.com/nathan-hello/htmx-template/src/auth"
	"github.com/nathan-hello/htmx-template/src/components"
	"github.com/nathan-hello/htmx-template/src/db"
	"github.com/nathan-hello/htmx-template/src/utils"
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
	for c := range m.clients {
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

// TODO: error handling to the client
func ChatSocket(w http.ResponseWriter, r *http.Request) {
	d := utils.Db()

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

		t := &utils.ChatMessage{}
		err = json.Unmarshal(clientMsg, &t)
		if err != nil {
			log.Println(err)
			return
		}
		if t.Text == "" {
			continue
		}
		if t.Author == "" {
			t.Author = "anon"
		}
		if t.Color == "" {
			t.Color = "bg-blue-200"
		}
		t.CreatedAt = time.Now()

		var buffMsg bytes.Buffer
		components.ChatMessage(t).Render(r.Context(), &buffMsg)
		log.Printf("buffMsg: %s, err: %#v\n", buffMsg.String(), err)

		manager.BroadcastMessage(buffMsg.Bytes())
		err = d.InsertMessage(r.Context(),
			db.InsertMessageParams{
				RoomID:    1,
				Author:    t.Author,
				Message:   t.Text,
				CreatedAt: t.CreatedAt,
			})
		if err != nil {
			log.Println(err)
			return
		}
	}

}

func Chat(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value(auth.ClaimsContextKey).(*auth.CustomClaims)
	if ok {
		w.Header().Set("HX-Redirect", "/")
		return
	}
	state := components.ClientState{
		IsAuthed: ok,
	}
	embed := r.URL.Query().Get("embed") == "true"

	if r.Method == "GET" {
		d := utils.Db()

		recents, err := d.SelectMessagesByChatroom(
			r.Context(),
			db.SelectMessagesByChatroomParams{
				RoomID: 1,
				Limit:  10,
			})
		if err != nil {
			log.Println(err)
		}

		var buffer bytes.Buffer
		for _, msg := range recents {
			components.ChatMessage(&utils.ChatMessage{
				Author:    msg.Author,
				Text:      msg.Message,
				Color:     "bg-blue-200",
				CreatedAt: msg.CreatedAt,
			}).Render(r.Context(), &buffer)
		}

		components.ChatRoot(state, embed).Render(r.Context(), w)
		w.Write(buffer.Bytes())

		return
	}
}
