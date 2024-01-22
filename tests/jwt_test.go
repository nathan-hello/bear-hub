package test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nathan-hello/htmx-template/src/db"
	"github.com/nathan-hello/htmx-template/src/utils"
)

func TestCreateAccess(t *testing.T) {
	access, err := utils.CreateAccess(uuid.New())

	if err != nil {
		t.Error(err)
	}

	t.Logf("TestCreateAccess/access: %#v\n", access)
}
func TestValidate(t *testing.T) {
	access, err := utils.CreateAccess(uuid.New())

	if err != nil {
		t.Error(err)
	}

	t.Logf("TestValidate/access: %#v\n", access)

	err = utils.VerifyJWT(access)

	if err != nil {
		t.Error(err)
	}

}

func TestTokenExpires(t *testing.T) {
	c := utils.Env()
	c.ACCESS_EXPIRY_TIME = time.Second * 1

	access, err := utils.CreateAccess(uuid.New())

	if err != nil {
		t.Error(err)
	}

	t.Logf("TestTokenExpires/access: %#v\n", access)

	err = utils.VerifyJWT(access)

	if err != nil {
		t.Error(err)
	}

	time.Sleep(time.Second * 2)

	err = utils.VerifyJWT(access)

	if err == nil {
		t.Error("Token is still valid even after waiting past expiration time")
		t.FailNow()
	}

	t.Logf("TestTokenExpires/msg: token successfully invalidated: %#v\n", err)
}

func TestDbJwt(t *testing.T) {
	ctx := context.Background()

	d, err := sql.Open("postgres", utils.Env().DB_URI)

	if err != nil {
		t.Error(err)
	}

	f := db.New(d)

	fullUser, err := f.InsertUser(ctx, db.InsertUserParams{
		Username:          "black-bear-test-1",
		EncryptedPassword: "honey",
	})

	defer func() {
		err = f.DeleteUser(ctx, fullUser.ID)

		if err != nil {
			t.Error(err)
		}
		fmt.Printf("deleted user: %#v\n", fullUser.ID.String())
	}()

	if err != nil {
		t.Error(err)
	}

	access, err := utils.CreateAccess(fullUser.ID)
	if err != nil {
		t.Error(err)
	}

	err = utils.VerifyJWT(access)
	if err != nil {
		t.Error(err)
	}

	err = f.InsertToken(ctx, db.InsertTokenParams{JwtType: "access", Jwt: access, Valid: true})
	if err != nil {
		t.Error(err)
	}

	token, err := f.SelectUsersTokens(ctx, fullUser.ID)
	if err != nil {
		t.Error(err)
	}

	if len(token) < 1 {
		t.Error("Token length 0: %#v\n", token)
	}

}
