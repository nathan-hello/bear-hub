// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0

package db

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Profile struct {
	Username sql.NullString
	Todos    sql.NullInt64
	ID       uuid.UUID
}

type Todo struct {
	ID        int64
	CreatedAt time.Time
	Body      string
}

type User struct {
	CreatedAt         time.Time
	Username          sql.NullString
	Email             sql.NullString
	EncryptedPassword string
	PasswordCreatedAt time.Time
	ID                uuid.UUID
}
