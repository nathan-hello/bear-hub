package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// TODO: all of these should be ErrBadLogin to prevent telling hostiles what is going on
var (
	ErrParsingJwt      = errors.New("internal Server Error - 11021")
	ErrInvalidToken    = errors.New("internal Server Error - 11022")
	ErrNoTokenInHeader = errors.New("internal Server Error - 11023")
	ErrBadJwtMethod    = errors.New("internal Server Error - 11025")
)

const (
	ErrJwtExpired       = "JWT Expired, %#v"
	ErrJwtSigningMethod = "unexpected signing method, %#v"
	ErrJwtNotValid      = "JWT not valid, %#v"
	ErrJwtParseFailed   = "JWT could not be parsed, %#v"
)

func CreateAccess(userId uuid.UUID) (string, error) {

	access := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(Env().ACCESS_EXPIRY_TIME).Unix(),
		"iat": time.Now().Unix(),
		"sub": userId,
		"iss": "no-magic-stack-example",
	})

	signedAccessToken, err := access.SignedString([]byte(Env().JWT_SECRET))

	if err != nil {
		return "", err
	}
	return signedAccessToken, nil
}

func CreateRefresh(userId uuid.UUID) (string, error) {

	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(Env().REFRESH_EXPIRY_TIME).Unix(),
		"iat": time.Now().Unix(),
		"sub": userId,
		"iss": "no-magic-stack-example",
	})

	signedRefreshToken, err := refresh.SignedString(Env().JWT_SECRET)

	if err != nil {
		return "", err
	}

	return signedRefreshToken, nil
}

func RefreshJwt(access string, refresh string) (string, error) {
	err := VerifyJWT(refresh)

	if err != nil {
		return "", err
	}

	return "", nil
}

func VerifyJWT(t string) error {

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
