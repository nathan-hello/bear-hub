package auth

import (
	"context"
	"database/sql"
	"slices"

	// "fmt"
	"net/mail"
	"time"

	"github.com/google/uuid"
	"github.com/nathan-hello/htmx-template/src/db"
	"github.com/nathan-hello/htmx-template/src/utils"
	"golang.org/x/crypto/bcrypt"
)

type AuthSignUp struct {
	Username    string
	UsernameErr string
	Email       string
	EmailErr    string
	Password    string
	PassErr     string
	PassConf    string
	PassConfErr string
	MiscErrs    []string
}

func (a *AuthSignUp) RenderErrs() []string {
	errs := []string{a.UsernameErr, a.EmailErr, a.PassErr, a.PassConfErr}
	errs = append(errs, a.MiscErrs...)
	for i, v := range errs {
		if v == "" {
			errs = slices.Delete(errs, i, i+1)
		}
	}
	return errs
}


func (a *AuthSignUp) FlushPasswords() {
	a.Password = ""
	a.PassConf = ""
}

func (a *AuthSignUp) validateStrings() bool {
	_, emailErr := mail.ParseAddress(a.Email)
	if emailErr != nil {
		if utils.AuthConfig.EmailRequired {
			a.EmailErr = utils.ErrEmailInvalid.Error()
		}
	}
	if len(a.Username) < 3 {
		if utils.AuthConfig.UsernameRequired {
			a.UsernameErr = utils.ErrUsernameTooShort.Error()

		}
	}

	if a.Username == "" && a.Email == "" {
		a.MiscErrs = append(a.MiscErrs, utils.ErrEmailOrUsernameReq.Error())
	}

	if len(a.Password) < 7 {
		a.PassErr = utils.ErrPasswordInvalid.Error()
	}
	if a.Password != a.PassConf {
		a.PassConfErr = utils.ErrPassNoMatch.Error()
	}

	ok := len(a.RenderErrs()) == 0
	if !ok {
		a.FlushPasswords()
		return false
	}

	return true
}

func (a *AuthSignUp) SignUp() *db.InsertUserRow {
	ok := a.validateStrings()
	if !ok {
		return nil
	}
	ctx := context.Background()

	if a.Email != "" {
		_, err := db.Db().SelectUserByEmail(ctx, a.Email)
		if err != sql.ErrNoRows {
			a.EmailErr = utils.ErrEmailTaken.Error()
		}
	}
	if a.Username != "" {
		_, err := db.Db().SelectUserByUsername(ctx, a.Username)
		if err != sql.ErrNoRows {
			a.UsernameErr = utils.ErrUsernameTaken.Error()
		}
	}

	userId := uuid.NewString()
	salt := uuid.NewString()[:8]
	pass, err := bcrypt.GenerateFromPassword([]byte(a.Password+salt), bcrypt.DefaultCost)

	if err != nil {
		a.MiscErrs = append(a.MiscErrs, utils.ErrHashPassword.Error())
		return nil
	}

	newUser, err := db.Db().InsertUser(
		ctx,
		db.InsertUserParams{
			ID:                userId,
			Email:             a.Email,
			Username:          a.Username,
			EncryptedPassword: string(pass),
			PasswordSalt:      salt,
			PasswordCreatedAt: time.Now(),
		})

	if err != nil {
		a.MiscErrs = append(a.MiscErrs, utils.ErrDbInsertUser.Error())
		return nil
	}

	return &newUser
}

type AuthSignIn struct {
	UserOrEmail    string
	UserOrEmailErr string
	Password       string
	PassErr        string
	MiscErrs       []string
}
func (a *AuthSignIn) RenderErrs() []string {
	errs := []string{a.UserOrEmailErr,a.PassErr}
	errs = append(errs, a.MiscErrs...)
	for i, v := range errs {
		if v == "" {
			errs = slices.Delete(errs, i, i+1)
		}
	}
	return errs
}
func (a *AuthSignIn) FlushPassword() {
	a.Password = ""
}

func (a *AuthSignIn) SignIn() *db.InsertUserRow {
	if a.UserOrEmail == "" || a.Password == "" {
		a.MiscErrs = append(a.MiscErrs, utils.ErrBadLogin.Error())
		a.FlushPassword()
		return nil
	}

	var user db.User
	ctx := context.Background()

	_, err := mail.ParseAddress(a.UserOrEmail)
	if err == nil {
		user, err = db.Db().SelectUserByEmailWithPassword(ctx, a.UserOrEmail)
		if err != nil {
			a.MiscErrs = append(a.MiscErrs, utils.ErrBadLogin.Error())
			a.FlushPassword()
			return nil
		}
	} else {
		user, err = db.Db().SelectUserByUsernameWithPassword(ctx, a.UserOrEmail)
		if err != nil {
			a.MiscErrs = append(a.MiscErrs, utils.ErrBadLogin.Error())
			a.FlushPassword()
			return nil
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte(a.Password+user.PasswordSalt))
	if err != nil {
		a.MiscErrs = append(a.MiscErrs, utils.ErrBadLogin.Error())
		a.FlushPassword()
		return nil
	}

	user.EncryptedPassword = ""

	return &db.InsertUserRow{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
	}
}
