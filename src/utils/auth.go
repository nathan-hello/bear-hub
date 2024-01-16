package utils

import (
	"context"
	"database/sql"
	"errors"

	// "fmt"
	"net/mail"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/nathan-hello/htmx-template/src/db"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUsernameTooShort = errors.New("Username too short")
	ErrPasswordTooShort = errors.New("Password too short")
	ErrEmailInvalid     = errors.New("Invalid email")
	ErrPasswordInvalid  = errors.New("Password invalid")
	ErrPassNoMatch      = errors.New("Passwords don't match")
	ErrBadLogin         = errors.New("Incorrect password or account does not exist")
	ErrDbConnection     = errors.New("Internal Server Error - 12482")
	ErrDbConnection2    = errors.New("Internal Server Error - 12483")
	ErrHashPassword     = errors.New("Internal Server Error - 19283")
	ErrHashPassword2    = errors.New("Internal Server Error - 19284")
	ErrDbInsertUser     = errors.New("Internal Server Error - 12382")
	ErrParseForm        = errors.New("Internal Server Error - 13481")
	ErrParseForm2       = errors.New("Internal Server Error - 13482")
)

const (
	FieldUsername = "username"
	FieldPassword = "password"
	FieldEmail    = "email"
	FieldPassConf = "password-confirmation"
	FieldUser     = "user"
)

var AllFields = []string{
	FieldUsername,
	FieldEmail,
	FieldPassConf,
	FieldPassword,
	FieldUser,
}

type AuthError struct {
	Field string
	Value string
	Err   error
}

type SignUpCredentials struct {
	Username string
	Password string
	PassConf string
	Email    string
}

type SignInCredentials struct {
	User string
	Pass string
}

func (c *SignUpCredentials) ValidateStrings() *[]AuthError {
	errs := []AuthError{}
	ok := true

	_, emailErr := mail.ParseAddress(c.Email)
	if c.Email != "" && emailErr != nil {
		errs = append(errs, AuthError{Field: FieldEmail, Err: ErrEmailInvalid, Value: c.Email})
		ok = false
	}

	if !(len(c.Username) > 3) {
		errs = append(errs, AuthError{Field: FieldUsername, Err: ErrUsernameTooShort, Value: c.Username})
		ok = false
	}

	if !(len(c.Password) > 7) {
		errs = append(errs, AuthError{Field: FieldPassword, Err: ErrPasswordTooShort, Value: ""})
		ok = false
	}

	if c.Password != c.PassConf {
		errs = append(errs, AuthError{Field: FieldPassConf, Err: ErrPassNoMatch, Value: ""})
		ok = false
	}

	if !ok {
		return &errs
	} else {
		return nil
	}
}

func (c *SignUpCredentials) SignUp() (string, *[]AuthError) {
	ctx := context.Background()
	errs := []AuthError{}

	d, err := sql.Open("postgres", Env().DB_URI)
	if err != nil {
		errs = append(errs, AuthError{Field: "", Err: ErrDbConnection, Value: ""})
		return "", &errs
	}

	conn := db.New(d)

	email := sql.NullString{String: c.Email, Valid: c.Email != ""}
	pass, err := bcrypt.GenerateFromPassword([]byte(c.Password), bcrypt.DefaultCost)
	if err != nil {
		errs = append(errs, AuthError{Field: "", Err: ErrHashPassword, Value: ""})
		return "", &errs
	}

	newUser, err := conn.InsertUser(
		ctx,
		db.InsertUserParams{
			Email:             email,
			Username:          c.Username,
			EncryptedPassword: string(pass),
			PasswordCreatedAt: time.Now(),
		})

	if err != nil {
		errs = append(errs, AuthError{Field: "", Err: ErrDbInsertUser, Value: ""})
		return "", &errs
	}

	return newUser.Username, nil

}

func (c *SignInCredentials) SignIn() (*db.User, *[]AuthError) {
	errs := []AuthError{}
	if c.User == "" || c.Pass == "" {
		errs = append(errs, AuthError{Field: FieldUser, Err: ErrBadLogin, Value: c.User})
		return nil, &errs
	}

	var user db.User
	ctx := context.Background()
	d, err := sql.Open("postgres", Env().DB_URI)
	if err != nil {
		errs = append(errs, AuthError{Err: ErrDbConnection2})
		return nil, &errs
	}

	conn := db.New(d)

	if _, err := mail.ParseAddress(c.User); err == nil {
		user, err = conn.SelectUserByEmail(ctx, sql.NullString{String: c.User, Valid: err == nil})
		if err != nil {
			errs = append(errs, AuthError{Field: FieldUser, Err: ErrBadLogin, Value: c.User})
			return nil, &errs
		}
	} else {
		user, err = conn.SelectUserByUsername(ctx, c.User)
		if err != nil {
			errs = append(errs, AuthError{Field: FieldUser, Err: ErrBadLogin, Value: c.User})
			return nil, &errs
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte(c.Pass))

	if err != nil {
		errs = append(errs, AuthError{Err: ErrHashPassword})
		return nil, &errs
	}

	user.EncryptedPassword = ""

	return &user, nil
}

const (
	ErrJwtExpired       = "JWT Expired, %#v"
	ErrJwtSigningMethod = "unexpected signing method, %#v"
	ErrJwtNotValid      = "JWT not valid, %#v"
	ErrJwtParseFailed   = "JWT could not be parsed, %#v"
)

func CreateAccess(userId uuid.UUID) (string, error) {

	access := jwt.NewWithClaims(jwt.SigningMethodEdDSA, jwt.MapClaims{
		"expires_at": time.Now().Add(time.Hour * 1).Unix(),
		"created_at": time.Now().Unix(),
		"user_id":    userId,
	})

	signedAccessToken, err := access.SignedString(Env().JWT_SECRET)

	if err != nil {
		return "", err
	}
	return signedAccessToken, nil
}

func CreateRefresh(userId uuid.UUID) (string, error) {

	refresh := jwt.NewWithClaims(jwt.SigningMethodEdDSA, jwt.MapClaims{
		"user_id":    userId,
		"expires_at": time.Now().Add(time.Hour * 72).Unix(),
		"created_at": time.Now().Unix(),
	})

	signedRefreshToken, err := refresh.SignedString(Env().JWT_SECRET)

	if err != nil {
		return "", err
	}

	return signedRefreshToken, nil
}

// func RefreshJwt(access string, refresh string) (string, error) {

// token, err := jwt.Parse(refresh, func(token *jwt.Token) (interface{}, error) {
// if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// return nil, fmt.Errorf(ErrJwtSigningMethod, token.Header["alg"])
// }
// return []byte(Env().JWT_SECRET), nil
// })

// if err != nil {
// return "", fmt.Errorf(ErrJwtParseFailed, token)
// }

// claims, ok := token.Claims.(jwt.MapClaims)

// if ok && token.Valid {
// if claims["expires_at"].(int64) < time.Now().Unix() {
// return "", fmt.Errorf(ErrJwtExpired, claims)
// }
// if !ok {
// return "", fmt.Errorf(ErrJwtNotValid, claims)
// }
// }

// newAccess, err := CreateAccess(claims["user_id"].(uuid.UUID))

// if err != nil {
// return nil, fmt.Errorf(err.Error())
// }

// return newPair, nil

// }
