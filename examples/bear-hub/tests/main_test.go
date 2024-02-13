package test

import (
	"context"
	"testing"

	"github.com/nathan-hello/htmx-template/examples/bear-hub/src/db"
	"github.com/nathan-hello/htmx-template/examples/bear-hub/src/utils"
)

func TestDatabaseConnection(t *testing.T) {
	ctx := context.Background()

	f, err := utils.Db()
	if err != nil {
		t.Error(err)
	}

	user, err := f.InsertUser(ctx, db.InsertUserParams{
		Username:          "black-bear-test-22121231321",
		EncryptedPassword: "honey",
	})

	defer func() {
		err = f.DeleteUser(ctx, user.ID)
		if err != nil {
			t.Error(err)
		}
	}()

	if err != nil {
		t.Error(err)
	}

	// fmt.Printf("New user: %#v\n", user)

	fullUser, err := f.SelectUserByUsername(ctx, user.Username)

	if err != nil {
		t.Error(err)
	}

	// fmt.Printf("full user %#v\n", fullUser)

	newTodo, err := f.InsertTodo(ctx, db.InsertTodoParams{Body: "eat honey", Username: fullUser.Username})

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

	_, err = f.SelectTodosByUsername(ctx, fullUser.Username)

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
