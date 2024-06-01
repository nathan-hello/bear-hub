package utils

import (
	"encoding/json"
	"fmt"
	"time"
)

const TimeFormat = time.RFC3339

type rawChatMessage struct {
	Author string `json:"msg-author"`
	Text   string `json:"msg-text"`
	Color  string `json:"msg-color"`
}

type ChatMessage struct {
	Raw       rawChatMessage
	Author    string    `json:"msg-author"`
	Text      string    `json:"msg-text"`
	Color     string    `json:"msg-color"`
	CreatedAt time.Time `json:"msg-time"`
}

func (c *ChatMessage) UnmarshalJSON(bits []byte) error {
	c.Raw = rawChatMessage{}
	err := json.Unmarshal(bits, &c.Raw)
	if err != nil {
		return err
	}

	if c.Raw.Text == "" {
		return ErrNoTextInChatMsg
	}

	if c.Raw.Author == "" {
		c.Raw.Author = "anon"
	}

	if c.Raw.Color == "" {
		c.Raw.Color = "text-gray-500"
	}

	c.Author = c.Raw.Author
	c.Text = c.Raw.Text
	c.Color = c.Raw.Color

	c.CreatedAt = time.Now().UTC()

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
