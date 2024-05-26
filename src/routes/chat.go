package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

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

		msg, err := utils.NewChatMsgFromBytes(clientMsg)
		if err != nil {
			log.Println(err)
			return
		}

		buffMsg := &bytes.Buffer{}

		components.ChatMessage(msg).Render(r.Context(), buffMsg) // write component to buffMsg
		manager.BroadcastMessage(buffMsg.Bytes())

	}
}

func ApiChat(w http.ResponseWriter, r *http.Request) {
	htmlResponse := r.Header.Get("Content-Type") == "text/html"
	jsonResponse := r.Header.Get("Content-Type") == "application/json"
	if !htmlResponse && !jsonResponse {
		htmlResponse = true
	}

	if r.Method == "POST" {
		c := utils.ChatMessage{}
		err := json.NewDecoder(r.Body).Decode(&c)
		if err != nil {
			w.Write([]byte(err.Error()))
		}
		err = c.Validate()
		if err != nil {
			w.Write([]byte(err.Error()))
		}

		var resp []byte

		// We need to send html to subscribers no matter what
		var htmlMsg bytes.Buffer
		components.ChatMessage(&c).Render(r.Context(), &htmlMsg)
		manager.BroadcastMessage(htmlMsg.Bytes())

		if htmlResponse {
			resp = htmlMsg.Bytes()
		}
		if jsonResponse {
			jason, err := json.Marshal(c)
			if err != nil {
				fmt.Fprintf(w, "{error: \"%v\"}", err)
			}
			resp = jason
		}
		w.Write(resp)
	}
}

func Chat(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value(auth.ClaimsContextKey).(*auth.CustomClaims)
	state := components.ClientState{
		IsAuthed: ok,
	}
	embed := r.URL.Query().Get("embed") == "true"

	if r.Method == "GET" {
		recents, err := db.Db().SelectMessagesByChatroom(
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
			m := &utils.ChatMessage{
				Author:    msg.Author,
				Text:      msg.Message,
				Color:     msg.Color,
				CreatedAt: utils.RenderTime(msg.CreatedAt),
			}
			err := m.Validate()
			if err != nil {
				continue
			}
			components.ChatMessage(m).Render(r.Context(), &buffer)
		}

		components.ChatRoot(state, embed).Render(r.Context(), w)
		w.Write(buffer.Bytes())

		return
	}
}
