package db

import (
	"context"

	"github.com/nathan-hello/htmx-template/src/db/output/sqlite"
)

type InsertChatroomParams struct {
	Name    string
	Creator string
}

type InsertChatroomReturn struct {
	Id int64
}

func (d *dbInsert) Chatroom(ctx context.Context, args InsertChatroomParams) (int64, error) {
	errs := DbError{}
	if d.slDb != nil {
		asdf, err := d.slDb.InsertChatroom(ctx, sqlite.InsertChatroomParams{Name: args.Name, Creator: args.Creator})
		errs.slErr = err

	}

}
