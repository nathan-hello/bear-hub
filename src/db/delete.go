package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func (d *dbDelete) Chatroom(ctx context.Context, id int64) error {
	errs := DbError{}
	if d.slDb != nil {
		errs.slErr = d.slDb.DeleteChatroom(ctx, id)
	}
	if d.pgDb != nil {
		errs.pgErr = d.pgDb.DeleteChatroom(ctx, id)
	}
	return errs.Unwrap()
}

func (d *dbDelete) Message(ctx context.Context, id int64) error {
	errs := DbError{}
	if d.slDb != nil {
		errs.slErr = d.slDb.DeleteMessage(ctx, id)
	}
	if d.pgDb != nil {
		errs.pgErr = d.pgDb.DeleteMessage(ctx, id)
	}
	return errs.Unwrap()
}

func (d *dbDelete) Todo(ctx context.Context, id int64) error {
	errs := DbError{}
	if d.slDb != nil {
		errs.slErr = d.slDb.DeleteTodo(ctx, id)
	}
	if d.pgDb != nil {
		errs.pgErr = d.pgDb.DeleteTodo(ctx, id)
	}
	return errs.Unwrap()
}

func (d *dbDelete) TokensByUserId(ctx context.Context, id uuid.UUID) error {
	errs := DbError{}
	if d.slDb != nil {
		return d.slDb.DeleteTokensByUserId(ctx, id.String())
	}
	if d.pgDb != nil {
		return d.pgDb.DeleteTokensByUserId(ctx, pgtype.UUID{Bytes: id, Valid: true})
	}
	return errs.Unwrap()
}

func (d *dbDelete) User(ctx context.Context, id uuid.UUID) error {
	errs := DbError{}
	if d.slDb != nil {
		errs.slErr = d.slDb.DeleteTokensByUserId(ctx, id.String())
	}
	if d.pgDb != nil {
		errs.pgErr = d.pgDb.DeleteTokensByUserId(ctx, pgtype.UUID{Bytes: id, Valid: true})
	}
	return errs.Unwrap()
}
