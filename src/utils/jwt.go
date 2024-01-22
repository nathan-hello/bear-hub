package utils

import (
	"context"
	"database/sql"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/nathan-hello/htmx-template/src/db"
)

type CustomClaims struct {
	jwt.RegisteredClaims
	UserId   uuid.UUID
	Username string
	Jwt_type string
}

func NewTokenPair(userId uuid.UUID, username string) (string, string, error) {
	ac := jwt.MapClaims{
		"exp":      time.Now().Add(Env().ACCESS_EXPIRY_TIME).Unix(),
		"iat":      time.Now().Unix(),
		"sub":      userId,
		"iss":      "no-magic-stack-example",
		"username": username,
	}

	access := jwt.NewWithClaims(jwt.SigningMethodHS256, &ac)

	as, err := access.SignedString([]byte(Env().JWT_SECRET))

	if err != nil {
		return "", "", err
	}

	rc := jwt.MapClaims{
		"exp":      time.Now().Add(Env().ACCESS_EXPIRY_TIME).Unix(),
		"iat":      time.Now().Unix(),
		"sub":      userId,
		"iss":      "no-magic-stack-example",
		"username": username,
	}

	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, &rc)
	rs, err := refresh.SignedString([]byte(Env().JWT_SECRET))
	if err != nil {
		return "", "", err
	}

	return as, rs, nil
}

func ParseToken(t string) (*CustomClaims, error) {
	token, err := jwt.Parse(
		t, func(token *jwt.Token) (interface{}, error) {
			ok := token.Method.Alg() == "HS256"
			if !ok {
				// this error will not show unless logged because
				// the jwt library wraps this error
				return nil, ErrBadJwtMethod
			}
			return []byte(Env().JWT_SECRET), nil
		})

	if err != nil {
		return nil, ErrParsingJwt2
	}

	if !token.Valid {
		return nil, ErrParsingJwt2
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, ErrParsingJwt2
	}

	return claims, nil

}

func InsertNewToken(t string, u uuid.UUID, jwt_type string) error {

	ctx := context.Background()

	d, err := sql.Open("postgres", Env().DB_URI)

	if err != nil {
		return ErrDbConnection3
	}

	f := db.New(d)

	token, err := f.InsertToken(ctx, db.InsertTokenParams{JwtType: jwt_type, Jwt: t, Valid: true})

	if err != nil {
		return ErrDbInsertToken
	}

	err = f.InsertUsersTokens(ctx, db.InsertUsersTokensParams{UserID: u, TokenID: sql.NullInt64{Int64: token, Valid: token > 0}})

	if err != nil {
		return ErrDbInsertUsersToken
	}

	return nil
}

func VerifyJWTLocal(t string) error {

	// QUESTION: does this function also validate the exp claim?
	// maybe because it's a standard jwt claim, it just knows to
	token, err := jwt.Parse(
		t, func(token *jwt.Token) (interface{}, error) {
			ok := token.Method.Alg() == "HS256"
			if !ok {
				// this error will not show unless logged because
				// the jwt library wraps this error
				return nil, ErrBadJwtMethod
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

func ValidJwtInDb(t string) error {

	ctx := context.Background()
	d, err := sql.Open("postgres", Env().DB_URI)
	if err != nil {
		return ErrDbConnection4
	}
	f := db.New(d)

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

func InvalidateTokenFromTokenString(t string) error {

	ctx := context.Background()
	d, err := sql.Open("postgres", Env().DB_URI)
	if err != nil {
		return ErrDbConnection5
	}
	f := db.New(d)
	token, err := f.SelectTokenFromJwtString(ctx, t)
	if err != nil {
		return ErrDbSelectJwt2
	}

	user, err := f.SelectUserIdFromToken(ctx, sql.NullInt64{Int64: token.ID, Valid: token.ID > 0})
	if err != nil {
		return ErrDbSelectUserFromToken
	}

	err = f.UpdateUserTokensToInvalid(ctx, user)
	if err != nil {
		return ErrDbUpdateTokensInvalid
	}

	return ErrJwtPairInvalid
}

func NewTokenPairFromRefreshString(r string) (string, string, error) {
	claims, err := ParseToken(r)
	if err != nil {
		return "", "", err
	}

	access, refresh, err := NewTokenPair(claims.UserId, claims.Username)
	if err != nil {
		return "", "", err
	}
	return access, refresh, nil

}

func ValidatePairOrRefresh(a string, r string) (string, string, error) {

	err := VerifyJWTLocal(a)
	if err == nil {
		err = VerifyJWTLocal(r)
		if err == nil {
			return a, r, nil
		}
		return "", "", ErrJwtGoodAccBadRef
	}

	err = VerifyJWTLocal(r)
	if err != nil {
		return "", "", err
	}

	err = ValidJwtInDb(r)
	if err != nil {
		if err == ErrJwtInvalidInDb {
			return "", "", InvalidateTokenFromTokenString(r)
		}
		return "", "", err
	}

	access, refresh, err := NewTokenPairFromRefreshString(r)
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}
