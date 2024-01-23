package routes

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
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

		access, refresh, err := utils.NewTokenPair(
			&utils.JwtParams{
				Username: username,
				UserId:   *userId,
				Family:   uuid.New(),
			})

		if err != nil {
			returnFormWithErrors(&[]utils.AuthError{
				{Err: err},
			})
		}

		utils.SetTokenCookies(w, access, refresh)
		w.Header().Set("HX-Redirect", fmt.Sprintf("/profile/%v", username))
		return

	}

	if r.Method == "GET" {
		if access, ok := utils.ValidateJwtOrDelete(w, r); ok {
			c, err := utils.ParseToken(access)
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
				{Err: utils.ErrParseForm},
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
		access, ok := utils.ValidateJwtOrDelete(w, r)
		if ok {
			p, err := utils.ParseToken(access)
			if err == nil {
				w.Header().Set("HX-Redirect", fmt.Sprintf("/profile/%v", p.Username))
				return
			}
			utils.DeleteJwtCookies(w)
		}
		components.SignInForm(nil).Render(r.Context(), w)
		return
	}
}

func SignOut(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		if r.URL.Path == "/signout" {
			w.Header().Set("HX-Redirect", "/?signedout=true")
		}
	}
}
