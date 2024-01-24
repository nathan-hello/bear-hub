package utils

import (
	"context"
	"database/sql"
	"fmt"

	// "fmt"
	"net/mail"
	"time"

	"github.com/google/uuid"
	"github.com/nathan-hello/htmx-template/src/db"
	"golang.org/x/crypto/bcrypt"
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

func (c *SignUpCredentials) validateStrings() *[]AuthError {
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
		fmt.Printf("pass: %#v\npassconf: %#v\n", c.Password, c.PassConf)
		errs = append(errs, AuthError{Field: FieldPassConf, Err: ErrPassNoMatch, Value: ""})
		ok = false
	}

	if !ok {
		return &errs
	} else {
		return nil
	}
}

func (c *SignUpCredentials) SignUp() (string, *uuid.UUID, *[]AuthError) {
	strErrs := c.validateStrings()
	if strErrs != nil {
		return "", nil, strErrs
	}
	ctx := context.Background()
	errs := []AuthError{}

	conn, err := Db()
	if err != nil {
		errs = append(errs, AuthError{Field: "", Err: err, Value: ""})
		return "", nil, &errs
	}

	email := sql.NullString{String: c.Email, Valid: c.Email != ""}
	pass, err := bcrypt.GenerateFromPassword([]byte(c.Password), bcrypt.DefaultCost)
	if err != nil {
		errs = append(errs, AuthError{Field: "", Err: ErrHashPassword, Value: ""})
		return "", nil, &errs
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
		return "", nil, &errs
	}

	return newUser.Username, &newUser.ID, nil

}

func (c *SignInCredentials) SignIn() (*db.User, *[]AuthError) {
	errs := []AuthError{}
	if c.User == "" || c.Pass == "" {
		errs = append(errs, AuthError{Field: FieldUser, Err: ErrBadLogin, Value: c.User})
		return nil, &errs
	}

	var user db.User
	ctx := context.Background()
	conn, err := Db()
	if err != nil {
		errs = append(errs, AuthError{Err: ErrDbConnection})
		return nil, &errs
	}

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
