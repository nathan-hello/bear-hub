// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0

package todos

import (
	"time"
)

type Todo struct {
	ID        int64
	CreatedAt time.Time
	Body      string
}
