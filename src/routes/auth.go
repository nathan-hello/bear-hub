package routes

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/nathan-hello/htmx-template/src/auth"
	"github.com/nathan-hello/htmx-template/src/components"
	"github.com/nathan-hello/htmx-template/src/db"
	"github.com/nathan-hello/htmx-template/src/utils"
)

func SignUp(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.ClaimsContextKey).(auth.CustomClaims)
	if ok {
		HandleRedirect(w, r, fmt.Sprintf("/profile/%s", claims.Username), nil)
		return
	}

	returnFormWithErrors := func(errs *[]auth.AuthError) {
		components.SignUpForm(components.RenderAuthError(errs)).Render(r.Context(), w)
		errs = nil
	}

	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			returnFormWithErrors(&[]auth.AuthError{
				{Err: utils.ErrParseForm},
			})
		}

		cred := auth.SignUpCredentials{
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

                jwtFamily := uuid.New()
		access, refresh, err := auth.NewTokenPair(
			&auth.JwtParams{
				Username: username,
				UserId:   *userId,
				Family:   jwtFamily,
			})
                utils.PrintlnOnDevMode("access, refresh", access, refresh)
		if err != nil {
			returnFormWithErrors(&[]auth.AuthError{
				{Err: err},
			})
		}

                d := utils.Db()
                _, err = d.InsertToken(r.Context(), db.InsertTokenParams{
                        JwtType: "access_token",
                        Jwt: access,
                        Valid: true,
                        Family: jwtFamily,
                })

		if err != nil {
			returnFormWithErrors(&[]auth.AuthError{
				{Err: err},
			})
		}

		auth.SetTokenCookies(w, access, refresh)
		w.Header().Set("HX-Redirect", "/")
		return

	}

	if r.Method == "GET" {
		components.SignUp().Render(r.Context(), w)
		return
	}
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.ClaimsContextKey).(auth.CustomClaims)
	if ok {
		HandleRedirect(w, r, fmt.Sprintf("/profile/%s", claims.Username), nil)
		return
	}

	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			components.SignInForm(components.RenderAuthError(&[]auth.AuthError{{Err: err}})).Render(r.Context(), w)
			return
		}

		cred := auth.SignInCredentials{
			User: r.FormValue("user"),
			Pass: r.FormValue("password"),
		}

		user, errs := cred.SignIn()

		if errs != nil {
			components.SignInForm(components.RenderAuthError(errs)).Render(r.Context(), w)
			return
		}

                jwtFamily := uuid.New()
		access, refresh, err := auth.NewTokenPair(
			&auth.JwtParams{
				Username: user.Username,
				UserId:   user.ID,
				Family:   jwtFamily,
			})

		if err != nil {
			components.SignInForm(components.RenderAuthError(errs)).Render(r.Context(), w)
			return
		}


		auth.SetTokenCookies(w, access, refresh)
		HandleRedirect(w, r, fmt.Sprintf("/profile/%s", claims.Username), nil)
		return
	}

	if r.Method == "GET" {
		components.SignIn().Render(r.Context(), w)
		return
	}
}

func SignOut(w http.ResponseWriter, r *http.Request) {
	auth.DeleteJwtCookies(w)
	HandleRedirect(w, r, "/", utils.ErrUserSignedOut)
}
