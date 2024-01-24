package utils

import "errors"

var (
	ErrUsernameTooShort = errors.New("username too short")
	ErrPasswordTooShort = errors.New("password too short")
	ErrEmailInvalid     = errors.New("invalid email")
	ErrPasswordInvalid  = errors.New("password invalid")
	ErrPassNoMatch      = errors.New("passwords don't match")
	ErrBadLogin         = errors.New("incorrect password or account does not exist")
	ErrParseForm        = errors.New("internal Server Error - 13481")
)

var (
	ErrHashPassword = errors.New("internal Server Error - 19283")
)

// TODO: all of these should be ErrBadLogin to prevent telling hostiles what is going on
var (
	ErrParsingJwt       = errors.New("internal Server Error - 11001")
	ErrInvalidToken     = errors.New("internal Server Error - 11002")
	ErrJwtNotInHeader   = errors.New("internal Server Error - 11003")
	ErrJwtNotInDb       = errors.New("internal Server Error - 11004")
	ErrJwtMethodBad     = errors.New("internal Server Error - 11005")
	ErrJwtInvalidInDb   = errors.New("internal Server Error - 11007")
	ErrJwtPairInvalid   = errors.New("internal Server Error - 11008")
	ErrJwtGoodAccBadRef = errors.New("internal Server Error - 11013")
)

var (
	ErrDbConnection          = errors.New("internal Server Error - 12401")
	ErrDbInsertUser          = errors.New("internal Server Error - 12405")
	ErrDbInsertToken         = errors.New("internal Server Error - 12406")
	ErrDbInsertProfile       = errors.New("internal Server Error - 12402")
	ErrDbInsertUsersToken    = errors.New("internal Server Error - 12407")
	ErrDbSelectJwt           = errors.New("internal Server Error - 12408")
	ErrDbInvalidateJwts      = errors.New("internal Server Error - 12409")
	ErrDbSelectUserFromToken = errors.New("internal Server Error - 12411")
	ErrDbUpdateTokensInvalid = errors.New("internal Server Error - 12412")
)
