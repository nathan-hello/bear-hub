package auth

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/nathan-hello/htmx-template/src/db"
	"github.com/nathan-hello/htmx-template/src/utils"
)

type JwtParams struct {
	Username string
	UserId   uuid.UUID
	Family   uuid.UUID
}

type ContextClaimType string

const ClaimsContextKey ContextClaimType = "claims"

type CustomClaims struct {
	jwt.RegisteredClaims
	UserId   uuid.UUID `json:"sub"`
	Username string    `json:"username"`
	Jwt_type string    `json:"jwt_type"`
	Family   uuid.UUID `json:"family"`
}

func NewTokenPair(j *JwtParams) (string, string, error) {
	ac := jwt.MapClaims{
		"exp":      time.Now().Add(utils.Env().ACCESS_EXPIRY_TIME).Unix(),
		"iat":      time.Now().Unix(),
		"iss":      "no-magic-stack-example",
		"sub":      j.UserId,
		"username": j.Username,
		"jwt_type": "access_token",
		"family":   j.Family,
	}

	access := jwt.NewWithClaims(jwt.SigningMethodHS256, &ac)

	as, err := access.SignedString([]byte(utils.Env().JWT_SECRET))
	if err != nil {
		return "", "", err
	}

	rc := jwt.MapClaims{
		"exp":      time.Now().Add(utils.Env().REFRESH_EXPIRY_TIME).Unix(),
		"iat":      time.Now().Unix(),
		"iss":      "no-magic-stack-example",
		"sub":      j.UserId,
		"username": j.Username,
		"jwt_type": "refresh_token",
		"family":   j.Family,
	}

	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, &rc)
	rs, err := refresh.SignedString([]byte(utils.Env().JWT_SECRET))
	if err != nil {
		return "", "", err
	}
	return as, rs, nil

}

func ParseToken(t string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(
		t,
		&CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			ok := token.Method.Alg() == "HS256"
			if !ok {
				// this error will not show unless logged because
				// the jwt library wraps this error
				return nil, utils.ErrJwtMethodBad
			}
			return []byte(utils.Env().JWT_SECRET), nil
		})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, utils.ErrParsingJwt
	}

	return claims, nil

}

func jwtValidLocal(t string) error {
	token, err := jwt.Parse(
		t, func(token *jwt.Token) (interface{}, error) {
			ok := token.Method.Alg() == "HS256"
			if !ok {
				// this error will not show unless logged because
				// the jwt library wraps this error
				return nil, utils.ErrJwtMethodBad
			}
			return []byte(utils.Env().JWT_SECRET), nil
		})

	if err != nil {
		return utils.ErrParsingJwt
	}

	if !token.Valid {
		return utils.ErrInvalidToken
	}

	return nil
}

func ValidateJwt(t string) error {
	err := jwtValidLocal(t)
	if err != nil {
		return err
	}

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

func InsertNewToken(t string, jwt_type string) error {
	claims, err := ParseToken(t)
	if err != nil {
		return err
	}
	ctx := context.Background()
	d := utils.Db()

	tokenId, err := d.InsertToken(ctx, db.InsertTokenParams{JwtType: jwt_type, Jwt: t, Valid: true, Family: claims.Family})
	if err != nil {
		return utils.ErrDbInsertToken
	}

	err = d.InsertUsersTokens(ctx, db.InsertUsersTokensParams{UserID: claims.UserId, TokenID: tokenId})
	if err != nil {
		return utils.ErrDbInsertUsersToken
	}
	return nil
}

func InvalidateJwtFamily(t string) error {

	ctx := context.Background()
	d := utils.Db()

	token, err := d.SelectTokenFromJwtString(ctx, t)
	if err != nil {
		return utils.ErrDbSelectJwt
	}

	claims, err := ParseToken(token.Jwt)
	if err != nil {
		return err
	}

	err = d.UpdateTokensFamilyInvalid(ctx, claims.Family)
	if err != nil {
		return utils.ErrDbUpdateTokensInvalid
	}

	return nil
}

func NewPairFromRefresh(r string) (string, string, error) {
	claims, err := ParseToken(r)
	if err != nil {
		return "", "", err
	}

	access, refresh, err := NewTokenPair(&JwtParams{UserId: claims.UserId, Username: claims.Username})
	if err != nil {
		return "", "", err
	}
	return access, refresh, nil

}

func ValidatePairOrRefresh(a string, r string) (string, string, error) {

	err := ValidateJwt(a)
        // if access is good, let's just refresh
	if err == nil {
		err = ValidateJwt(r)
                // if refresh and access are good, return to sender
		if err == nil {
			return a, r, nil
		}
                // if access is good but refresh is bad, we don't refresh based off
                // of access tokens, so it's better to just error and reauth
                utils.PrintlnOnDevMode("errjwtgoodaccbadref:", err)
		return "", "", utils.ErrJwtGoodAccBadRef
	}

        // even if access was bad, maybe the refresh is good
	err = ValidateJwt(r)
	if err != nil {
		if err == utils.ErrJwtInvalidInDb {
			return "", "", InvalidateJwtFamily(r)
		}
		return "", "", err
	}

        // sweet, a good refresh jwt. let's make a new pair
	access, refresh, err := NewPairFromRefresh(r)
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}
