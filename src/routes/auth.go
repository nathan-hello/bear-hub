package routes

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/nathan-hello/htmx-template/src/components"
	"github.com/nathan-hello/htmx-template/src/db"
	"github.com/nathan-hello/htmx-template/src/utils"
)

func SignUp(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/signup" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	ctx := context.Background()

	returnFormWithErrors := func(errs *utils.AuthErrors) {
		response, err := templ.ToGoHTML(ctx, components.SignUpForm(*errs))

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write([]byte(response))
		return
	}

	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		cred := utils.SignUpCredentials{
			Username: r.FormValue("username"),
			Password: r.FormValue("password"),
			PassConf: r.FormValue("password-confirm"),
			Email:    r.FormValue("email"),
		}

		errs := cred.ValidateStrings()

		if len(errs.ErrsStr) > 0 {
			returnFormWithErrors(errs)
			return
		}

		errs = cred.ValidateDatabase()

		if len(errs.ErrsStr) > 0 {
			returnFormWithErrors(errs)
			return
		}


		w.Header().Set("HX-Redirect", fmt.Sprintf("/profile/%v", newUser.Username))
		return

	}

	if r.Method == "GET" {

		defaultFormErr := utils.FormErr{BorderColor: "border-blue-500"}
		form := utils.AuthErrors{
			Email:    defaultFormErr,
			Username: defaultFormErr,
			Password: defaultFormErr,
			PassConf: defaultFormErr,
		}

		response, err := templ.ToGoHTML(ctx, components.SignUp(form))

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write([]byte(response))
		return

	}
}

func SignIn(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/signin" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	ctx := context.Background()
	d, err := sql.Open("postgres", utils.Env().DB_URI)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	conn := db.New(d)

	returnFormWithErrors := func(errs *utils.AuthErrors) {
		response, err := templ.ToGoHTML(ctx, components.SignInForm(errs))

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write([]byte(response))
		return
	}

	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		cred := utils.SignInCredentials{
			EmailOrUsername: r.FormValue("EmailOrUsername"),
			Password:        r.FormValue("password"),
		}

		errs := cred.ValidateDatabase()

		if len(errs.ErrsStr) > 0 {
			returnFormWithErrors(errs)
			return
		}

		isEmail := cred.IsEmail()

		var user db.User

		user := conn.SelectUserWithEmailPassword()

		w.Header().Set("HX-Redirect", fmt.Sprintf("/profile/%v", cred.Username))
		return

	}

	if r.Method == "GET" {

		defaultFormErr := utils.FormErr{BorderColor: "border-blue-500"}
		errs := utils.SignInErrors{
			Email:    defaultFormErr,
			Password: defaultFormErr,
		}

		response, err := templ.ToGoHTML(ctx, components.SignIn(&errs))

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write([]byte(response))
		return

	}
}
