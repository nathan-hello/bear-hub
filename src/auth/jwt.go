package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/nathan-hello/htmx-template/src/utils"
)

type JwtParams struct {
	Username string
	UserId   string
	Family   string
}

func NewTokenPair(j *JwtParams) (string, string, error) {
	if j.Family == "" {
		j.Family = uuid.New().String()
	}
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
	err = DbInsertNewToken(as, "access")
	if err != nil {
		return "", "", err
	}
	err = DbInsertNewToken(rs, "refresh")
	if err != nil {
		return "", "", err
	}
	return as, rs, nil

}

func ParseToken(t string) (*utils.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(
		t,
		&utils.CustomClaims{},
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

	claims, ok := token.Claims.(*utils.CustomClaims)
	if !ok {
		return nil, utils.ErrParsingJwt
	}

	return claims, nil

}

func localValidateJwt(t string) error {
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

func ValidateJwtFromString(t string) error {
	err := localValidateJwt(t)
	if err != nil {
		return err
	}
	err = DbValidateJwt(t)
	return err
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

	err := ValidateJwtFromString(a)
	// if access is good, let's just refresh
	if err == nil {
		err = ValidateJwtFromString(r)
		// if refresh and access are good, return to sender
		if err == nil {
			return a, r, nil
		}
		// if access is good but refresh is bad, we don't refresh based off
		// of access tokens, so it's better to just error and reauth
		return "", "", utils.ErrJwtGoodAccBadRef
	}

	// even if access was bad, maybe the refresh is good
	err = ValidateJwtFromString(r)
	if err != nil {
		if err == utils.ErrJwtInvalidInDb {
			return "", "", DbInvalidateJwtFamily(r)
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
