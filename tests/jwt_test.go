package test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nathan-hello/htmx-template/src/auth"
	"github.com/nathan-hello/htmx-template/src/db"
	"github.com/nathan-hello/htmx-template/src/utils"
)

func TestNewPairAndParse(t *testing.T) {
	initEnv(t)
	access, refresh, err := auth.NewTokenPair(
		&auth.JwtParams{
			UserId:   uuid.New().String(),
			Username: "black-bear",
			Family:   uuid.New().String(),
		})

	if err != nil {
		t.Error(err)
	}

	// t.Logf("TestNewPairAndParse/access: %v\n", access)
	// t.Logf("TestNewPairAndParse/refresh: %v\n", refresh)

	_, err = auth.ParseToken(access)
	if err != nil {
		t.Error(err)
	}

	_, err = auth.ParseToken(refresh)
	if err != nil {
		t.Error(err)
	}

	// t.Logf("TestNewPairAndParse/ap: %v\n", ap)
	// t.Logf("TestNewPairAndParse/rp: %v\n", rp)

}
func TestJwtExpiry(t *testing.T) {
	initEnv(t)
	c := utils.Env()
	c.ACCESS_EXPIRY_TIME = time.Second * 1
	c.REFRESH_EXPIRY_TIME = time.Second * 2

	access, refresh, err := auth.NewTokenPair(
		&auth.JwtParams{
			UserId:   uuid.New().String(),
			Username: "black-bear",
			Family:   uuid.New().String(),
		})

	if err != nil {
		t.Error(err)
	}

	// t.Logf("access: %#v\n", access)
	// t.Logf("refresh: %#v\n", refresh)

	_, err = auth.ParseToken(access)

	if err != nil {
		t.Error(err)
	}

	time.Sleep(time.Millisecond * 1100)

	_, err = auth.ParseToken(access)

	if err == nil {
		t.Error("access is still valid even after waiting past expiration time")
	}

	// t.Logf("token successfully invalidated: %#v\n", err)

	_, err = auth.ParseToken(refresh)

	if err != nil {
		t.Error(err)
	}

	time.Sleep(time.Second * 1)

	_, err = auth.ParseToken(refresh)

	if err == nil {
		t.Error("refresh is still valid even after waiting past expiration time")
	}

	// t.Logf("refresh successfully invalidated: %#v\n", err)
}

func TestDbJwt(t *testing.T) {
	initEnv(t)
	ctx := context.Background()

	fullUser, err := db.Db().InsertUser(ctx, db.InsertUserParams{
		Username:          "black-bear-test-1",
		EncryptedPassword: "honey",
	})

	defer func() {
		err = db.Db().DeleteUser(ctx, fullUser.ID)

		if err != nil {
			t.Error(err)
		}
		// fmt.Printf("deleted user: %#v\n", fullUser.ID.String())
	}()

	if err != nil {
		t.Error(err)
	}

	access, refresh, err := auth.NewTokenPair(
		&auth.JwtParams{
			Username: fullUser.Username,
			UserId:   fullUser.ID,
		})
	if err != nil {
		t.Error(err)
	}

	_, err = auth.ParseToken(access)
	if err != nil {
		t.Error(err)
	}
	_, err = auth.ParseToken(refresh)
	if err != nil {
		t.Error(err)
	}

	err = auth.DbInsertNewToken(access, "access_token")
	defer func() {
		err := db.Db().DeleteTokensByUserId(ctx, fullUser.ID)
		if err != nil {
			t.Error(err)
		}
	}()
	if err != nil {
		t.Error(err)
	}
	err = auth.DbInsertNewToken(refresh, "refresh_token")
	if err != nil {
		t.Error(err)
	}

	tokens, err := db.Db().SelectUsersTokens(ctx, fullUser.ID)
	if err != nil {
		t.Error(err)
	}

	if len(tokens) == 0 {
		t.Error("Token length 0\n")
	}

	err = db.Db().UpdateUserTokensToInvalid(ctx, fullUser.ID)
	if err != nil {
		t.Error(err)
	}

	_, err = db.Db().SelectUsersTokens(ctx, fullUser.ID)
	if err != nil {
		t.Error(err)
	}

}
