package routes

import (
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nathan-hello/htmx-template/src/utils"
)

var (
	ErrParsingJwt      = errors.New("Internal Server Error - 11021\n")
	ErrInvalidToken    = errors.New("Internal Server Error - 11022\n")
	ErrNoTokenInHeader = utils.ErrBadLogin
	ErrBadJwtMethod    = errors.New("Internal Server Error - 11024\n")
)

func verifyJWT(request *http.Request) error {

	if _, ok := request.Header["Token"]; !ok {
		return ErrNoTokenInHeader
	}

	token, err := jwt.Parse(
		request.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodECDSA)
			if !ok {
				return nil, ErrBadJwtMethod
			}
			return []byte(utils.Env().JWT_SECRET), nil
		})

	if err != nil {
		return ErrParsingJwt
	}

	if !token.Valid {
		return ErrInvalidToken
	}

	return nil

}
