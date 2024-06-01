package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	gws "github.com/gorilla/websocket"
	"github.com/nathan-hello/htmx-template/src/components"
	"github.com/nathan-hello/htmx-template/src/db"
	"github.com/nathan-hello/htmx-template/src/utils"
)

const DEFAULT_ROOM_ID = 1
const DEFAULT_CHAT_COLOR = "text-gray-500"

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
	state := utils.GetClientState(r)
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

		var msg utils.ChatMessage
		err = json.Unmarshal(clientMsg, &msg)
		if err != nil {
			log.Println(err)
			w.Write([]byte(err.Error()))
		}

		db.Db().InsertMessage(
			r.Context(),
			db.InsertMessageParams{
				AuthorID:       &state.UserId,
				AuthorUsername: msg.Author,
				Message:        msg.Text,
				CreatedAt:      msg.CreatedAt,
				RoomID:         DEFAULT_ROOM_ID,
			})

		buffMsg := &bytes.Buffer{}

		components.ChatMessage(&msg).Render(r.Context(), buffMsg) // write component to buffMsg
		manager.BroadcastMessage(buffMsg.Bytes())

	}
}

func ApiChat(w http.ResponseWriter, r *http.Request) {
	//state := utils.GetClientState(r)
	htmlResponse := r.Header.Get("Content-Type") == "text/html"
	jsonResponse := r.Header.Get("Content-Type") == "application/json"
	if !htmlResponse && !jsonResponse {
		htmlResponse = true
	}

	if r.Method == "POST" {
		c := utils.ChatMessage{}
		err := json.NewDecoder(r.Body).Decode(&c)
		if err != nil {
			fmt.Fprintf(w, "{error: \"%v\"}\n", err)
			return
		}

		var resp []byte

		// We need to send html to subscribers no matter what
		htmlMsg := bytes.Buffer{}
		components.ChatMessage(&c).Render(r.Context(), &htmlMsg)
		manager.BroadcastMessage(htmlMsg.Bytes())

		db.Db().InsertMessage(
			r.Context(),
			db.InsertMessageParams{
				AuthorID:       nil,
				AuthorUsername: c.Author,
				Message:        c.Text,
				RoomID:         DEFAULT_ROOM_ID,
				CreatedAt:      c.CreatedAt,
			})

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
	state := utils.GetClientState(r)
	embed := r.URL.Query().Get("embed") == "true"

	if r.Method == "GET" {
		recents, err := db.Db().SelectMessagesByChatroom(
			r.Context(),
			db.SelectMessagesByChatroomParams{
				RoomID: DEFAULT_ROOM_ID,
				Limit:  10,
			})
		if err != nil {
			log.Println(err)
		}

		var buffer bytes.Buffer
		for _, msg := range recents {
			var color string
			if msg.ChatroomColor == nil {
				color = DEFAULT_CHAT_COLOR
			} else {
				color = *msg.ChatroomColor
			}
			m := &utils.ChatMessage{
				Author:    msg.AuthorUsername,
				Text:      msg.Message,
				Color:     color,
				CreatedAt: msg.CreatedAt,
			}
			components.ChatMessage(m).Render(r.Context(), &buffer)
		}

		components.ChatRoot(state, embed).Render(r.Context(), w)
		w.Write(buffer.Bytes())

		return
	}
}
