package routes

import (
	"context"
	"fmt"
	"net/http"

	"github.com/nathan-hello/htmx-template/src/auth"
	"github.com/nathan-hello/htmx-template/src/components"
	"github.com/nathan-hello/htmx-template/src/utils"
)

func signUpErrMsg(err error, errs *[]auth.AuthError, ctx context.Context, w http.ResponseWriter) {
	if err != nil {
		components.SignUpForm(components.RenderAuthError(&[]auth.AuthError{{Err: err}})).Render(ctx, w)
	}
	if errs != nil {
		components.SignUpForm(components.RenderAuthError(errs)).Render(ctx, w)
	}
}

func signInErrMsg(err error, errs *[]auth.AuthError, ctx context.Context, w http.ResponseWriter) {
	if err != nil {
		components.SignInForm(components.RenderAuthError(&[]auth.AuthError{{Err: err}})).Render(ctx, w)
	}
	if errs != nil {
		components.SignInForm(components.RenderAuthError(errs)).Render(ctx, w)
	}
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.ClaimsContextKey).(auth.CustomClaims)
	if ok {
		HandleRedirect(w, r, fmt.Sprintf("/profile/%s", claims.Username), nil)
		return
	}
	state := components.ClientState{
		IsAuthed: ok,
	}

	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			signUpErrMsg(utils.ErrParseForm, nil, r.Context(), w)
		}

		cred := auth.SignUpCredentials{
			Username: r.FormValue("username"),
			Password: r.FormValue("password"),
			PassConf: r.FormValue("password-confirmation"),
			Email:    r.FormValue("email"),
		}

		username, userId, errs := cred.SignUp()

		if errs != nil {
			signUpErrMsg(nil, errs, r.Context(), w)
			return
		}

		access, refresh, err := auth.NewTokenPair(
			&auth.JwtParams{
				Username: username,
				UserId:   userId.String(),
			})

		if err != nil {
			signUpErrMsg(err, nil, r.Context(), w)
			return
		}

		auth.SetTokenCookies(w, access, refresh)
		w.Header().Set("HX-Redirect", "/")
		return

	}

	if r.Method == "GET" {
		components.SignUp(state).Render(r.Context(), w)
		return
	}
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value(auth.ClaimsContextKey).(*auth.CustomClaims)
	if ok {
		w.Header().Set("HX-Redirect", "/")
		return
	}
	state := components.ClientState{
		IsAuthed: ok,
	}

	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			signInErrMsg(err, nil, r.Context(), w)
			return
		}

		cred := auth.SignInCredentials{
			User: r.FormValue("user"),
			Pass: r.FormValue("password"),
		}

		user, errs := cred.SignIn()

		if errs != nil {
			signInErrMsg(nil, errs, r.Context(), w)
			return
		}

		access, refresh, err := auth.NewTokenPair(
			&auth.JwtParams{
				Username: user.Username,
				UserId:   user.ID,
			})

		if err != nil {
			signInErrMsg(err, nil, r.Context(), w)
			return
		}

		claims, err := auth.ParseToken(access)
		if err != nil {
			signInErrMsg(err, nil, r.Context(), w)
			return
		}

		auth.SetTokenCookies(w, access, refresh)
		HandleRedirect(w, r, fmt.Sprintf("/profile/%s", claims.Username), nil)
		return
	}

	if r.Method == "GET" {
		components.SignIn(state).Render(r.Context(), w)
		return
	}
}

func SignOut(w http.ResponseWriter, r *http.Request) {
	auth.DeleteJwtCookies(w)
	HandleRedirect(w, r, "/", utils.ErrUserSignedOut)
}
