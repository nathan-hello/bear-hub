package routes

import (
	"fmt"
	"net/http"

	"github.com/nathan-hello/htmx-template/src/components"
	"github.com/nathan-hello/htmx-template/src/utils"
)

func SignUp(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/signup" {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("HX-Redirect", "404")
		return
	}

	returnFormWithErrors := func(errs *[]utils.AuthError) {
		components.SignUpForm(components.RenderAuthError(errs)).Render(r.Context(), w)
	}

	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			returnFormWithErrors(&[]utils.AuthError{
				{Err: utils.ErrParseForm},
			})
		}

		cred := utils.SignUpCredentials{
			Username: r.FormValue("username"),
			Password: r.FormValue("password"),
			PassConf: r.FormValue("password-confirm"),
			Email:    r.FormValue("email"),
		}

		errs := cred.ValidateStrings()

		if errs != nil {
			returnFormWithErrors(errs)
			return
		}

		username, errs := cred.SignUp()

		if errs != nil {
			returnFormWithErrors(errs)
			return
		}

		w.Header().Set("HX-Redirect", fmt.Sprintf("/profile/%v", username))
		return

	}

	if r.Method == "GET" {
		components.SignUpForm(nil).Render(r.Context(), w)
	}
}

func SignIn(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/signin" {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("HX-Redirect", "404")
		return
	}

	returnFormWithErrors := func(errs *[]utils.AuthError) {
		components.SignInForm(components.RenderAuthError(errs)).Render(r.Context(), w)
	}

	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			returnFormWithErrors(&[]utils.AuthError{
				{Err: utils.ErrParseForm2},
			})
			return
		}

		cred := utils.SignInCredentials{
			User: r.FormValue("user"),
			Pass: r.FormValue("password"),
		}

		username, errs := cred.SignIn()

		if errs != nil {
			returnFormWithErrors(errs)
			return
		}

		w.Header().Set("HX-Redirect", fmt.Sprintf("/profile/%v", username))
		return

	}

	if r.Method == "GET" {
		components.SignInForm(nil).Render(r.Context(), w)
	}
}
