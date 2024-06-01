package utils

import (
	"encoding/json"
	"fmt"
	"time"
)

const TimeFormat = time.RFC3339

type rawChatMessage struct {
	Author    string `json:"msg-author"`
	Text      string `json:"msg-text"`
	Color     string `json:"msg-color"`
}

type ChatMessage struct {
        raw  rawChatMessage
	Author    string `json:"msg-author"`
	Text      string `json:"msg-text"`
	Color     string `json:"msg-color"`
	CreatedAt *time.Time `json:"msg-time"`
}

func (c *ChatMessage) UnmarshalJSON(bits []byte) error {
        var raw rawChatMessage
	err := json.Unmarshal(bits, &raw)
	if err != nil {
		return err
	}
        if raw.Text == "" {
                return ErrNoTextInChatMsg
        }
        
        c.Author = raw.Author
        c.Text = raw.Text
        c.Color = raw.Color

        if c.Author == "" {
                c.Author = "anon"
        }

        if c.Color == "" {
                c.Color = "text-gray-500"
        }
        
        utcTime  := time.Now().UTC()
        c.CreatedAt = &utcTime

	return nil
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
