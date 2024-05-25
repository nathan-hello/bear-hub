package test

import (
	"testing"

	"github.com/nathan-hello/htmx-template/src/db"
	"github.com/nathan-hello/htmx-template/src/utils"
)

func initEnv(t *testing.T) {
	err := utils.InitEnv("../.env")
	if err != nil {
		t.Fatal(err)
	}
	_ = utils.Env()

	err = db.DbInit()
	if err != nil {
		t.Fatal(err)
	}

	if db.Db() == nil {
		t.Fatal("Db() returned a nil *db.Queries")
	}
}
