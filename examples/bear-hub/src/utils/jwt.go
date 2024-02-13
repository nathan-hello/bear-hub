package utils

import (
	"context"
	"database/sql"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/nathan-hello/htmx-template/examples/bear-hub/examples/bear-hub/src/db"
)

type JwtParams struct {
	Username string
	UserId   uuid.UUID
	Family   uuid.UUID
}

type CustomClaims struct {
	jwt.RegisteredClaims
	UserId   uuid.UUID `json:"sub"`
	Username string    `json:"username"`
	Jwt_type string    `json:"jwt_type"`
	Family   uuid.UUID `json:"family"`
}

func NewTokenPair(j *JwtParams) (string, string, error) {
	ac := jwt.MapClaims{
		"exp":      time.Now().Add(Env().ACCESS_EXPIRY_TIME).Unix(),
		"iat":      time.Now().Unix(),
		"iss":      "no-magic-stack-example",
		"sub":      j.UserId,
		"username": j.Username,
		"jwt_type": "access_token",
		"family":   j.Family,
	}

	access := jwt.NewWithClaims(jwt.SigningMethodHS256, &ac)

	as, err := access.SignedString([]byte(Env().JWT_SECRET))
	if err != nil {
		return "", "", err
	}

	rc := jwt.MapClaims{
		"exp":      time.Now().Add(Env().REFRESH_EXPIRY_TIME).Unix(),
		"iat":      time.Now().Unix(),
		"iss":      "no-magic-stack-example",
		"sub":      j.UserId,
		"username": j.Username,
		"jwt_type": "refresh_token",
		"family":   j.Family,
	}

	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, &rc)
	rs, err := refresh.SignedString([]byte(Env().JWT_SECRET))
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
				return nil, ErrJwtMethodBad
			}
			return []byte(Env().JWT_SECRET), nil
		})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, ErrParsingJwt
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
				return nil, ErrJwtMethodBad
			}
			return []byte(Env().JWT_SECRET), nil
		})

	if err != nil {
		return ErrParsingJwt
	}

	if !token.Valid {
		return ErrInvalidToken
	}

	return nil
}

func ValidateJwt(t string) error {
	err := jwtValidLocal(t)
	if err != nil {
		return err
	}

	ctx := context.Background()
	f, err := Db()
	if err != nil {
		return ErrDbConnection
	}

	token, err := f.SelectTokenFromJwtString(ctx, t)

	if err != nil {
		if err == sql.ErrNoRows {
			return ErrJwtNotInDb
		}
		return ErrDbSelectJwt
	}
	if !token.Valid {
		return ErrJwtInvalidInDb
	}

	return nil
}

func InsertNewToken(t string, jwt_type string) error {
	claims, err := ParseToken(t)
	if err != nil {
		return err
	}
	ctx := context.Background()
	f, err := Db()
	if err != nil {
		return ErrDbConnection
	}

	tokenId, err := f.InsertToken(ctx, db.InsertTokenParams{JwtType: jwt_type, Jwt: t, Valid: true, Family: claims.Family})
	if err != nil {
		return ErrDbInsertToken
	}

	err = f.InsertUsersTokens(ctx, db.InsertUsersTokensParams{UserID: claims.UserId, TokenID: tokenId})
	if err != nil {
		return ErrDbInsertUsersToken
	}
	return nil
}

func InvalidateJwtFamily(t string) error {

	ctx := context.Background()
	f, err := Db()
	if err != nil {
		return ErrDbConnection
	}

	token, err := f.SelectTokenFromJwtString(ctx, t)
	if err != nil {
		return ErrDbSelectJwt
	}

	claims, err := ParseToken(token.Jwt)
	if err != nil {
		return err
	}

	err = f.UpdateTokensFamilyInvalid(ctx, claims.Family)
	if err != nil {
		return ErrDbUpdateTokensInvalid
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
	if err == nil {
		err = ValidateJwt(r)
		if err == nil {
			return a, r, nil
		}
		return "", "", ErrJwtGoodAccBadRef
	}

	err = ValidateJwt(r)
	if err != nil {
		if err == ErrJwtInvalidInDb {
			return "", "", InvalidateJwtFamily(r)
		}
		return "", "", err
	}

	access, refresh, err := NewPairFromRefresh(r)
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}
