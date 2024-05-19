package utils

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/nathan-hello/htmx-template/src/db"
)

func NewChatMsg(bits []byte) (*ChatMessage, error) {
        t := &ChatMessage{}
        err := json.Unmarshal(bits, &t)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		if t.Text == "" {
			return nil, ErrNoTextInChatMsg
		}
		if t.Author == "" {
			t.Author = "anon"
		}
		if t.Color == "" {
			t.Color = "bg-blue-200"
		}
                t.InsertTimeNow()
        return t, nil
}

type ChatMessage struct {
	Author    string `json:"msg-author"`
	Text      string `json:"msg-text"`
	Color     string `json:"msg-color"`
        CreatedAt string `json:"msg-time"`
}

        

func (c *ChatMessage) InsertTimeNow() {
        c.CreatedAt = time.Now().UTC().Format(time.RFC3339)
}

func (c *ChatMessage) InsertDbAsync(ctx context.Context) {

	d := Db()
        err := d.InsertMessage(ctx,
			db.InsertMessageParams{
				RoomID:    1,
				Author:    c.Author,
				Message:   c.Text,
			})
        if err != nil {
                log.Println(err)
        }
}
