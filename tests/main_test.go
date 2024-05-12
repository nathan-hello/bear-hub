package test

import (
	"context"
	"testing"

	"github.com/nathan-hello/htmx-template/src/db"
	"github.com/nathan-hello/htmx-template/src/utils"
)

func TestDatabaseConnection(t *testing.T) {
	initEnv(t)
	ctx := context.Background()
	d := utils.Db()

	user, err := d.InsertUser(ctx, db.InsertUserParams{
		Username:          "black-bear-test-22121231321",
		EncryptedPassword: "honey",
	})

	defer func() {
		err = d.DeleteUser(ctx, user.ID)
		if err != nil {
			t.Error(err)
		}
	}()

	if err != nil {
		t.Error(err)
	}

	// fmt.Printf("New user: %#v\n", user)

	fullUser, err := d.SelectUserByUsername(ctx, user.Username)

	if err != nil {
		t.Error(err)
	}

	// fmt.Printf("full user %#v\n", fullUser)

	newTodo, err := d.InsertTodo(ctx, db.InsertTodoParams{Body: "eat honey", Username: fullUser.Username})

	defer func() {
		err = d.DeleteTodo(ctx, newTodo.ID)

		if err != nil {
			t.Error(err)
		}
	}()

	if err != nil {
		t.Error(err)
	}

	// fmt.Printf("newTodo: %#v\n", newTodo)

	_, err = d.SelectTodosByUsername(ctx, fullUser.Username)

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
