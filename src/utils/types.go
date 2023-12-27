package utils

import (
	"context"
	"database/sql"
	"fmt"
	"net/mail"

	"github.com/nathan-hello/htmx-template/src/sqlc"
)

// Form submission validation belongs in routes/auth.go
// It's in here to prevent circular imports between components and routes.

type FormErr struct {
	Value       string
	BorderColor string
}

type SignUpErrors struct {
	Email    FormErr
	Username FormErr
	Password FormErr
	PassConf FormErr
	ErrsStr  []string
}

type Credentials struct {
	Username string
	Password string
	PassConf string
	Email    string
}

func (c *Credentials) Validate() *SignUpErrors {
	errs := SignUpErrors{
		Email:    FormErr{BorderColor: "bg-blue-500", Value: c.Email},
		Username: FormErr{BorderColor: "bg-blue-500", Value: c.Username},
		Password: FormErr{BorderColor: "bg-blue-500", Value: ""},
		PassConf: FormErr{BorderColor: "bg-blue-500", Value: ""},
	}

	user := len(c.Username) > 3
	pass := len(c.Password) > 4
	_, err := mail.ParseAddress(c.Email)

	if !user {
		errs.Username = FormErr{
			Value:       "",
			BorderColor: "border-red-500",
		}
		errs.ErrsStr = append(errs.ErrsStr, "Username is too short")

	}
	if !pass {
		errs.Password = FormErr{
			Value:       "",
			BorderColor: "border-red-500",
		}
		errs.ErrsStr = append(errs.ErrsStr, "Password is too short")
	}
	if err != nil {
		errs.Email = FormErr{
			Value:       "",
			BorderColor: "border-red-500",
		}
		errs.ErrsStr = append(errs.ErrsStr, "Email invalid")
	}
	if c.Password != c.PassConf {
		errs.PassConf = FormErr{
			Value:       "",
			BorderColor: "border-red-500",
		}
		errs.ErrsStr = append(errs.ErrsStr, "Passwords don't match")
	}

	return &errs
}

func (c *Credentials) ValidateEmailUsername() *SignUpErrors {
	errs := SignUpErrors{
		Email:    FormErr{BorderColor: "bg-blue-500", Value: c.Email},
		Username: FormErr{BorderColor: "bg-blue-500", Value: c.Username},
		Password: FormErr{BorderColor: "bg-blue-500", Value: ""},
		PassConf: FormErr{BorderColor: "bg-blue-500", Value: ""},
	}

	ctx := context.Background()
	db, err := sql.Open("postgres", Env().DB_URI)

	if err != nil {
		errs.ErrsStr = append(errs.ErrsStr, "Internal Server Error - 135232")
		return &errs
	}

	conn := sqlc.New(db)
	_, err = conn.SelectEmailAlreadyExists(ctx, sql.NullString{String: c.Email, Valid: c.Email != ""})
	emailOk := false

	if err != nil {
		if err == sql.ErrNoRows {
			emailOk = true
		} else {
			errs.ErrsStr = append(errs.ErrsStr, "Internal Server Error - 135234")
			return &errs
		}
	}

	if !emailOk {
		errs.ErrsStr = append(errs.ErrsStr, fmt.Sprintf("Email %v already has an account", c.Email))
		return &errs
	}

	usernameAlreadyExists, err := conn.SelectProfileByUsername(ctx, c.Username)

	if err != nil && err != sql.ErrNoRows {
		errs.ErrsStr = append(errs.ErrsStr, "Internal Server Error - 135236")
		return &errs
	}

	if usernameAlreadyExists.Username != "" {
		errs.ErrsStr = append(errs.ErrsStr, "Username is taken")
		return &errs
	}

	return &errs

}

func (c *Credentials) CustomErrorMessage(err string) *SignUpErrors {

	errs := SignUpErrors{
		Email:    FormErr{BorderColor: "bg-blue-500", Value: c.Email},
		Username: FormErr{BorderColor: "bg-blue-500", Value: c.Username},
		Password: FormErr{BorderColor: "bg-blue-500", Value: ""},
		PassConf: FormErr{BorderColor: "bg-blue-500", Value: ""},
	}
	errs.ErrsStr = append(errs.ErrsStr, err)
	return &errs

}
