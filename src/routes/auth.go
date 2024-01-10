package routes

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/google/uuid"
	"github.com/nathan-hello/htmx-template/src/components"
	"github.com/nathan-hello/htmx-template/src/db"
	"github.com/nathan-hello/htmx-template/src/utils"
	"github.com/nedpals/supabase-go"
)

func SignUp(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/signup" {
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

	returnFormWithErrors := func(errs *utils.SignUpErrors) {
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

		url := utils.Env().SUPABASE_POSTGREST_URL
		publicKey := utils.Env().SUPABASE_PUBLIC_KEY
		client := supabase.CreateClient(url, publicKey)

		user, err := client.Auth.SignUp(ctx, supabase.UserCredentials{
			Email:    cred.Email,
			Password: cred.Password,
		})

		uid, err := uuid.Parse(user.ID)

		if err != nil {
			returnFormWithErrors(cred.CustomErrorMessage(fmt.Sprintf("Internal Server Error - 482012 %#v", err)))
			return
		}

		uname, err := conn.UpdateProfileUsername(ctx, db.UpdateProfileUsernameParams{Username: cred.Username, ID: uid})

		if err != nil {
			returnFormWithErrors(cred.CustomErrorMessage("Internal Server Error - 125632"))
			return
		}

		w.Header().Set("HX-Redirect", fmt.Sprintf("/profile/%v", uname))
		return

	}

	if r.Method == "GET" {

		defaultFormErr := utils.FormErr{BorderColor: "border-blue-500"}
		form := utils.SignUpErrors{
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

	returnFormWithErrors := func(errs *utils.SignInErrors) {
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
			Email:    r.FormValue("email"),
			Password: r.FormValue("password"),
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

		url := utils.Env().SUPABASE_POSTGREST_URL
		publicKey := utils.Env().SUPABASE_PUBLIC_KEY
		client := supabase.CreateClient(url, publicKey)

		user, err := client.Auth.SignIn(ctx, supabase.UserCredentials{
			Email:    cred.Email,
			Password: cred.Password,
		})
		if err != nil {
			// returnFormWithErrors(cred.CustomErrorMessage("Incorrect password or account does not exist"))
			returnFormWithErrors(cred.CustomErrorMessage(fmt.Sprintf("%#v\n", err)))
			return
		}

		userUuid, err := uuid.Parse(user.User.ID)

		if err != nil {
			returnFormWithErrors(cred.CustomErrorMessage("Internal Server Error - 481234"))
			return
		}

		profileId, err := conn.SelectProfileById(ctx, userUuid)

		if err != nil {
			// returnFormWithErrors(cred.CustomErrorMessage("Incorrect password or account does not exist"))
			returnFormWithErrors(cred.CustomErrorMessage(fmt.Sprintf("%#v\n", err)))
			return
		}

		w.Header().Set("HX-Redirect", fmt.Sprintf("/profile/%v", profileId))
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
