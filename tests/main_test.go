package test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/lib/pq"
	"github.com/nathan-hello/htmx-template/src/db"
	"github.com/nathan-hello/htmx-template/src/utils"
)

func TestDatabaseConnection(t *testing.T) {
	ctx := context.Background()

	d, err := sql.Open("postgres", utils.Env().DB_URI)

	if err != nil {
		t.Error(err)
	}

	f := db.New(d)

	user, err := f.InsertUser(ctx, db.InsertUserParams{
		Username:          "black-bear-test-2",
		EncryptedPassword: "honey",
	})

	if err != nil {
		t.Error(err)
	}

	// fmt.Printf("New user: %#v\n", user)

	fullUser, err := f.SelectUserByUsername(ctx, user.Username)

	if err != nil {
		t.Error(err)
	}

	// fmt.Printf("full user %#v\n", fullUser)

	newProfile, err := f.InsertProfile(ctx, fullUser.ID)

	if err != nil {
		t.Error(err)
	}

	defer func() {
		err := f.DeleteProfile(ctx, newProfile)

		if err != nil {
			t.Error(err)
		}
		fmt.Printf("deleted profile: %#v\n", newProfile.String())

		err = f.DeleteUser(ctx, fullUser.ID)

		if err != nil {
			t.Error(err)
		}
		fmt.Printf("deleted user: %#v\n", fullUser.ID.String())
	}()

	// fmt.Printf("newProfile: %#v\n", newProfile)

	newTodo, err := f.InsertTodo(ctx, db.InsertTodoParams{Body: "eat honey", Author: fullUser.ID})

	defer func() {
		err = f.DeleteTodo(ctx, newTodo.ID)

		if err != nil {
			t.Error(err)
		}
	}()

	if err != nil {
		t.Error(err)
	}

	// fmt.Printf("newTodo: %#v\n", newTodo)

	_, err = f.SelectUserTodos(ctx, fullUser.ID)

	if err != nil {
		t.Error(err)
	}

	// fmt.Printf("printing up to 10 todos\n")
	//
	// for i, v := range rows {
	// 	fmt.Printf("row %v: %#v", i, v)
	// 	if i >= 10 {
	// 		break
	// 	}
	// }

}
