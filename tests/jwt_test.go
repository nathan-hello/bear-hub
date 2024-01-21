package test

import (
	"testing"
	"time"

	"github.com/google/uuid"
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
