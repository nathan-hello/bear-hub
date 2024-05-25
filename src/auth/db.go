package auth

import (
	"context"

	"github.com/nathan-hello/htmx-template/src/db"
	"github.com/nathan-hello/htmx-template/src/utils"
)

func DbInsertNewToken(t string, jwt_type string) error {
	claims, err := ParseToken(t)
	if err != nil {
		return err
	}
	ctx := context.Background()

	token, err := db.Db().InsertToken(ctx, db.InsertTokenParams{JwtType: jwt_type, Jwt: t, Valid: true, Family: claims.Family})
	if err != nil {
		return utils.ErrDbInsertToken
	}

	err = db.Db().InsertUsersTokens(ctx, db.InsertUsersTokensParams{UserID: claims.UserId, TokenID: token.ID})
	if err != nil {
		return utils.ErrDbInsertUsersToken
	}
	return nil
}

func DbValidateJwt(t string) error {
	// TODO: validate jwt in database
	// ctx := context.Background()
	// d := Db()
	//
	// token, err := d.SelectTokenFromJwtString(ctx, t)
	//
	// if err != nil {
	// 	if err == sql.ErrNoRows {
	// 		return ErrJwtNotInDb
	// 	}
	// 	return ErrDbSelectJwt
	// }
	// if !token.Valid {
	// 	return ErrJwtInvalidInDb
	// }

	return nil
}

func DbInvalidateJwtFamily(t string) error {

	ctx := context.Background()

	token, err := db.Db().SelectTokenFromJwtString(ctx, t)
	if err != nil {
		return utils.ErrDbSelectJwt
	}

	claims, err := ParseToken(token.Jwt)
	if err != nil {
		return err
	}
	err = db.Db().UpdateTokensFamilyInvalid(ctx, claims.Family)
	if err != nil {
		return utils.ErrDbUpdateTokensInvalid
	}

	return nil
}
