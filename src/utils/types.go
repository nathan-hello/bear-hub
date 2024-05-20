package utils

import (
	"encoding/json"
	"time"
)

type ChatMessage struct {
	Author    string `json:"msg-author"`
	Text      string `json:"msg-text"`
	Color     string `json:"msg-color"`
	CreatedAt string `json:"msg-time"`
}

func NewChatMsgFromBytes(bits []byte) (*ChatMessage, error) {
	t := &ChatMessage{}
	err := json.Unmarshal(bits, &t)
	if err != nil {
		return nil, err
	}
	err = t.Validate()
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (c *ChatMessage) Validate() error {
	if c.Text == "" {
		return ErrNoTextInChatMsg
	}
	if c.Author == "" {
		c.Author = "anon"
	}
	if c.Color == "" {
		c.Color = "bg-blue-200"
	}
	if c.CreatedAt == "" {
		c.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	}
	return nil

}
