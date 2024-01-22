package routes

import (
	"fmt"
	"net/http"

	"github.com/nathan-hello/htmx-template/src/components"
	"github.com/nathan-hello/htmx-template/src/utils"
)

func SignUp(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/signup" {
		redirectNotFound(w, r)
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

		username, userId, errs := cred.SignUp()

		if errs != nil {
			returnFormWithErrors(errs)
			return
		}

		access, refresh, err := utils.NewTokenPair(*userId, username)
		if err != nil {
			returnFormWithErrors(&[]utils.AuthError{
				{Err: err},
			})
		}

		SetTokenCookies(w, access, refresh)
		w.Header().Set("HX-Redirect", fmt.Sprintf("/profile/%v", username))
		return

	}

	if r.Method == "GET" {
		if ValidateOrRefreshPairFromCookies(w, r) {

			access, err := r.Cookie("access_token")
			if err != nil {
				w.Header().Set("HX-Redirect", "500")
			}

			c, err := utils.ParseToken(access.Value)
			if err != nil {
				w.Header().Set("HX-Redirect", "500")
			}

			w.Header().Set("HX-Redirect", fmt.Sprintf("/profile/%v", c.Username))
		}
		components.SignUpForm(nil).Render(r.Context(), w)
	}
}

func SignIn(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/signin" {
		redirectNotFound(w, r)
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
		if ValidateOrRefreshPairFromCookies(w, r) {

			access, err := r.Cookie("access_token")
			if err != nil {
				w.Header().Set("HX-Redirect", "500")
			}

			c, err := utils.ParseToken(access.Value)
			if err != nil {
				w.Header().Set("HX-Redirect", "500")
			}

			w.Header().Set("HX-Redirect", fmt.Sprintf("/profile/%v", c.Username))
		}
		components.SignInForm(nil).Render(r.Context(), w)
	}
}
