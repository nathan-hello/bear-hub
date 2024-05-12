package utils

import "time"

type ChatMessage struct {
	Author    string `json:"msg-author"`
	Text      string `json:"msg-text"`
	Color     string `json:"msg-color"`
	CreatedAt time.Time
}
