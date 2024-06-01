package utils

import (
	"encoding/json"
	"fmt"
	"time"
)

const TimeFormat = time.RFC3339

type rawChatMessage struct {
	Text   string `json:"msg-text"`
}

type ChatMessage struct {
        UserId string `json:"msg-userid"`
        Username string `json:"msg-author"`
	Text      string    `json:"msg-text"`
	Color     string    `json:"msg-color"`
	CreatedAt time.Time `json:"msg-time"`
}

func NewChatFromBytes(bits []byte, username string, userId string, color string) (*ChatMessage, error) {
        var raw rawChatMessage
        err := json.Unmarshal(bits, &raw)
        if err != nil {
                return nil, err
        }
        if raw.Text == "" {
                return nil, ErrNoTextInChatMsg
        }

        if username == "" {
                username = "anon"
        }

        fmt.Println(color)

        return &ChatMessage{
                UserId: userId,
                Username: username,
                Text: raw.Text,
                Color: color,
                CreatedAt: time.Now().UTC(),
        }, nil
}

// TimeToString(true)  = string(HH:MM)
//
// TimeToString(false) = .Format(time.RFC3339)
func (c *ChatMessage) TimeToString(hourMinute bool) string {
	if hourMinute {
		return fmt.Sprintf("%02d:%02d", c.CreatedAt.Hour(), c.CreatedAt.Minute())
	}
	return c.CreatedAt.Format(time.RFC3339)
}
