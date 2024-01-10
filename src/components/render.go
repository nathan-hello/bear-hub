package components

import (
	"time"
)

func formatTime(t *time.Time) string {
	return t.Format(time.RFC822)
}
