package auth

import (
	"context"
	"fmt"

	// "fmt"
	"net/mail"
	"time"

	"github.com/google/uuid"
	"github.com/nathan-hello/htmx-template/src/db"
	"github.com/nathan-hello/htmx-template/src/utils"
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
		errs = append(errs, AuthError{Field: FieldEmail, Err: utils.ErrEmailInvalid, Value: c.Email})
		ok = false
	}

	if len(c.Username) < 3 {
		errs = append(errs, AuthError{Field: FieldUsername, Err: utils.ErrUsernameTooShort, Value: c.Username})
		ok = false
	}

	if len(c.Password) < 7 {
		errs = append(errs, AuthError{Field: FieldPassword, Err: utils.ErrPasswordTooShort, Value: ""})
		ok = false
	}

	if c.Password != c.PassConf {
		fmt.Printf("pass: %#v\npassconf: %#v\n", c.Password, c.PassConf)
		errs = append(errs, AuthError{Field: FieldPassConf, Err: utils.ErrPassNoMatch, Value: ""})
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

	pass, err := bcrypt.GenerateFromPassword([]byte(c.Password), bcrypt.DefaultCost)
	if err != nil {
		errs = append(errs, AuthError{Field: "", Err: utils.ErrHashPassword, Value: ""})
		return "", nil, &errs
	}

	newUser, err := db.Db().InsertUser(
		ctx,
		db.InsertUserParams{
			Email:             c.Email,
			Username:          c.Username,
			EncryptedPassword: string(pass),
			PasswordCreatedAt: time.Now(),
		})

	if err != nil {
		errs = append(errs, AuthError{Field: "", Err: utils.ErrDbInsertUser, Value: ""})
		return "", nil, &errs
	}

	parsedId, err := uuid.Parse(newUser.ID)
	if err != nil {
		utils.PrintlnOnDevMode("insert user", err)

		errs = append(errs, AuthError{Field: "", Err: utils.ErrDbInsertUser, Value: ""})
		return "", nil, &errs
	}

	return newUser.Username, &parsedId, nil

}

func (c *SignInCredentials) SignIn() (*db.User, *[]AuthError) {
	errs := []AuthError{}
	if c.User == "" || c.Pass == "" {
		errs = append(errs, AuthError{Field: FieldUser, Err: utils.ErrBadLogin, Value: c.User})
		return nil, &errs
	}

	var user db.User
	ctx := context.Background()

	if _, err := mail.ParseAddress(c.User); err == nil {
		user, err = db.Db().SelectUserByEmail(ctx, c.User)
		if err != nil {
			errs = append(errs, AuthError{Field: FieldUser, Err: utils.ErrBadLogin, Value: c.User})
			return nil, &errs
		}
	} else {
		user, err = db.Db().SelectUserByUsername(ctx, c.User)
		if err != nil {
			errs = append(errs, AuthError{Field: FieldUser, Err: utils.ErrBadLogin, Value: c.User})
			return nil, &errs
		}
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte(c.Pass))

	if err != nil {
		errs = append(errs, AuthError{Err: utils.ErrHashPassword})
		return nil, &errs
	}

	user.EncryptedPassword = ""

	return &user, nil
}
