package utils

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/mail"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/nathan-hello/htmx-template/src/db"
	"golang.org/x/crypto/bcrypt"
)

// Form submission validation belongs in routes/auth.go
// It's in here to prevent circular imports between components and routes.

var (
	ErrUsernameTooShort = errors.New("Username too short")
	ErrPasswordTooShort = errors.New("Password too short")
	ErrEmailInvalid     = errors.New("Invalid email")
	ErrPasswordInvalid  = errors.New("Password invalid")
	ErrPassNoMatch      = errors.New("Passwords don't match")
	ErrBadLogin         = errors.New("Incorrect password or account does not exist")
	ErrDbConnection     = errors.New("Internal Server Error - 12482")
	ErrHashPassword     = errors.New("Internal Server Error - 19283")
	ErrDbInsertUser     = errors.New("Internal Server Error - 12382")
)

const (
	FieldUsername = "username"
	FieldPassword = "password"
	FieldEmail    = "email"
	FieldPassConf = "password-confirmation"
)

type AuthError struct {
	Field string
	Value string
	Err   error
}

type AuthCredentials struct {
	Username string
	Password string
	PassConf string
	Email    string
}

func (c *AuthCredentials) SignInStrings() *[]AuthError {

	user := len(c.Username) > 3
	pass := len(c.Password) > 7
	_, emailErr := mail.ParseAddress(c.Email)
	errs := []AuthError{}
	ok := true
	if !user {
		errs = append(errs, AuthError{Field: FieldUsername, Err: ErrUsernameTooShort, Value: c.Username})
		ok = false
	}
	if !pass {
		errs.Error(SectionPassword, ErrPasswordTooShort, "")
	}
	if c.Email != "" && emailErr != nil {
		errs.Error(SectionEmail, ErrEmailInvalid, c.Email)
	}
	if c.Password != c.PassConf {
		errs.Error(SectionPassConf, ErrPassNoMatch, "")
	}

	if !ok {
		return &errs
	} else {
		return nil
	}
}

func (c *SignUpCredentials) ValidateDatabase() *AuthErrors {
	errs := AuthErrors{}
	ok := true
	ctx := context.Background()
	d, err := sql.Open("postgres", Env().DB_URI)

	if err != nil {
		errs.Error("", ErrDbConnection, "")
		ok = false
		return &errs
	}

	conn := db.New(d)
	_, err = conn.SelectEmailOrUsernameAlreadyExists(ctx,
		db.SelectEmailOrUsernameAlreadyExistsParams{
			Email:    sql.NullString{String: c.Email, Valid: c.Email != ""},
			Username: c.Username,
		})

	if err != nil && err != sql.ErrNoRows {
		errs.Error("", ErrDbConnection, "")
		ok = false
		return &errs
	}

	if !ok {
		return &errs
	} else {
		return nil
	}

}

func (c *SignUpCredentials) MiscErrorMessage(err string) *AuthErrors {
	return &AuthErrors{ErrsStr: []string{err}}
}

func (c *SignUpCredentials) SignUp() (string, *AuthErrors) {
	ctx := context.Background()
	d, err := sql.Open("postgres", Env().DB_URI)

	errs := AuthErrors{}
	if err != nil {
		errs.Error("", ErrDbConnection, "")
		return "", &errs
	}

	conn := db.New(d)

	email := sql.NullString{String: c.Email, Valid: c.Email != ""}
	pass, err := bcrypt.GenerateFromPassword([]byte(c.Password), bcrypt.DefaultCost)
	if err != nil {
		errs.Error("", ErrHashPassword, "")
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
		errs.Error("", ErrDbInsertUser, "")
		return "", &errs
	}

	return newUser.Username, &errs

}

func (c *SignInCredentials) IsEmail() bool {
	_, err := mail.ParseAddress(c.EmailOrUsername)
	return err == nil
}

func (c *SignInCredentials) SignIn() (string, *AuthErrors) {
	var user db.User
	errs := AuthErrors{}

	ctx := context.Background()
	d, err := sql.Open("postgres", Env().DB_URI)

	if err != nil {
		errs.Error("", ErrBadLogin, "")
	}

	conn := db.New(d)

	if c.IsEmail() {
		userDb, err := conn.SelectUserByEmail(ctx, sql.NullString{String: c.EmailOrUsername, Valid: c.IsEmail()})
		user = userDb
		if err != nil {
			errs.Error("", ErrBadLogin, "")
		}
	} else {
		userDb, err := conn.SelectUserByUsername(ctx, c.EmailOrUsername)
		user = userDb
		if err != nil {
			errs.Error("", ErrBadLogin, "")
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte(c.Password))

	if err != nil {
		errs.Error("", ErrBadLogin, "")
	}

	return user.Username, &errs

}

type JwtStrings struct {
	access  string
	refresh string
}

const (
	ErrJwtExpired       = "JWT Expired, %#v"
	ErrJwtSigningMethod = "unexpected signing method, %#v"
	ErrJwtNotValid      = "JWT not valid, %#v"
	ErrJwtParseFailed   = "JWT could not be parsed, %#v"
)

func CreateJwt(userId uuid.UUID) (*JwtStrings, error) {

	access := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"expires_at": time.Now().Add(time.Hour * 1).Unix(),
		"created_at": time.Now().Unix(),
		"user_id":    userId,
	})

	signedAccessToken, err := access.SignedString(Env().JWT_SECRET)

	if err != nil {
		return nil, err
	}

	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    userId,
		"expires_at": time.Now().Add(time.Hour * 72).Unix(),
		"created_at": time.Now().Unix(),
	})

	signedRefreshToken, err := refresh.SignedString(Env().JWT_SECRET)

	if err != nil {
		return nil, err
	}

	return &JwtStrings{access: signedAccessToken, refresh: signedRefreshToken}, nil
}

func RefreshJwt(j JwtStrings) (*JwtStrings, error) {

	token, err := jwt.Parse(j.refresh, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf(ErrJwtSigningMethod, token.Header["alg"])
		}
		return []byte(Env().JWT_SECRET), nil
	})

	if err != nil {
		return nil, fmt.Errorf(ErrJwtParseFailed, token)
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		if claims["expires_at"].(int64) < time.Now().Unix() {
			return nil, fmt.Errorf(ErrJwtExpired, claims)
		}
		if !ok {
			return nil, fmt.Errorf(ErrJwtNotValid, claims)
		}
	}

	newPair, err := CreateJwt(claims["user_id"].(uuid.UUID))

	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	return newPair, nil

}
