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
		Redirect404(w, r)
		return
	}

	returnFormWithErrors := func(errs *[]utils.AuthError) {
		fmt.Printf("%#v\n", *errs)
		components.SignUpForm(components.RenderAuthError(errs)).Render(r.Context(), w)
		errs = nil
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
			PassConf: r.FormValue("password-confirmation"),
			Email:    r.FormValue("email"),
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
				Redirect401(w, r)
				return
			}

			w.Header().Set("HX-Redirect", fmt.Sprintf("/profile/%v", c.Username))
		}
		components.SignUp().Render(r.Context(), w)
	}
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	returnFormWithErrors := func(errs *[]utils.AuthError) {
		components.SignInForm(components.RenderAuthError(errs)).Render(r.Context(), w)
	}

	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			returnFormWithErrors(&[]utils.AuthError{
				{Err: utils.ErrBadLogin},
			})
			return
		}

		cred := utils.SignInCredentials{
			User: r.FormValue("user"),
			Pass: r.FormValue("password"),
		}

		user, errs := cred.SignIn()

		if errs != nil {
			returnFormWithErrors(errs)
			return
		}

		access, refresh, err := utils.NewTokenPair(
			&utils.JwtParams{
				Username: user.Username,
				UserId:   user.ID,
				Family:   uuid.New(),
			})

		if err != nil {
			returnFormWithErrors(&[]utils.AuthError{
				{Err: utils.ErrBadLogin},
			})
		}

		utils.SetTokenCookies(w, access, refresh)
		w.Header().Set("HX-Redirect", fmt.Sprintf("/profile/%v", user.Username))
		return
	}

	if r.Method == "GET" {
		components.SignIn().Render(r.Context(), w)
		return
	}
}

func SignOut(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		if r.URL.Path == "/signout" {
			RedirectSignOut(w, r)
		}
	}
}
