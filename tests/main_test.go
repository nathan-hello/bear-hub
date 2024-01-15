package test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/lib/pq"
	sqlc "github.com/nathan-hello/htmx-template/src/db"
	"github.com/nathan-hello/htmx-template/src/utils"
)

func TestDatabaseConnection(t *testing.T) {
	ctx := context.Background()

	db, err := sql.Open("postgres", utils.Env().DB_URI)

	if err != nil {
		panic(err)
	}

	f := sqlc.New(db)

	user, err := f.InsertUser(ctx, sqlc.InsertUserParams{
		Username:          "black-bear-test-2",
		EncryptedPassword: "honey",
	})

	defer func() {
		err = f.DeleteUser(ctx, user.ID)

		if err != nil {
			panic(err)
		}
	}()

	if err != nil {
		panic(err)
	}

	fmt.Printf("New user: %#v\n", user)

	fullUser, err := f.SelectUserByUsername(ctx, user.Username)

	if err != nil {
		panic(err)
	}

	fmt.Printf("full user %#v\n", fullUser)

	newProfile, err := f.InsertProfile(ctx, fullUser.ID)

	defer func() {
		err = f.DeleteProfile(ctx, newProfile)
		if err != nil {
			panic(err)
		}
	}()

	fmt.Printf("newProfile: %#v\n", newProfile)

	newTodo, err := f.InsertTodo(ctx, sqlc.InsertTodoParams{Body: "eat honey", Author: fullUser.ID})

	defer func() {
		err = f.DeleteTodo(ctx, newTodo.ID)

		if err != nil {
			panic(err)
		}
	}()

	if err != nil {
		panic(err)
	}

	fmt.Printf("newTodo: %#v\n", newTodo)

	rows, err := f.SelectAllTodos(ctx)

	if err != nil {
		panic(err)
	}

	fmt.Printf("printing up to 10 todos\n")

	for i, v := range rows {
		fmt.Printf("row %v: %#v", i, v)
		if i >= 10 {
			break
		}
	}

}
