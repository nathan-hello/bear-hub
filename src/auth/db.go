package auth

import (
	"context"
	"database/sql"

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

func DbValidateJwt(t string, user string) error {
	ctx := context.Background()
	
	token, err := db.Db().SelectTokenFromJwtString(ctx, t)

	if err != nil {
		if err == sql.ErrNoRows {
			return utils.ErrJwtNotInDb
		}
		return utils.ErrDbSelectJwt
	}
	if !token.Valid {
		return utils.ErrJwtInvalidInDb
	}

        _, err = db.Db().SelectUserById(context.Background(), user)

        if err != nil {
                return err
        }

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
