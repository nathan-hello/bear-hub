package test

import (
	"testing"

	"github.com/nathan-hello/htmx-template/src/utils"
)

func initEnv(t *testing.T) {
	err := utils.InitEnv("../.env")
	if err != nil {
		t.Fatal(err)
	}
	_ = utils.Env()

	err = utils.DbInit()
	if err != nil {
		t.Fatal(err)
	}

	conn := utils.Db()
	if conn == nil {
		t.Fatal("Db() returned a nil *db.Queries")
	}
}
