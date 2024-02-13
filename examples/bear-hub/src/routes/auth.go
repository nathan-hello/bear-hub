package routes

import (
	"fmt"
	"net/http"

	"bear-hub/src/components"
	"bear-hub/src/utils"

	"github.com/google/uuid"
)

func SignUp(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(utils.ClaimsContextKey).(utils.CustomClaims)
	if ok {
		HandleRedirect(w, r, fmt.Sprintf("/profile/%s", claims.Username), http.StatusSeeOther, nil)
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
		components.SignUp().Render(r.Context(), w)
		return
	}
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(utils.ClaimsContextKey).(utils.CustomClaims)
	if ok {
		HandleRedirect(w, r, fmt.Sprintf("/profile/%s", claims.Username), http.StatusSeeOther, nil)
		return
	}

	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			components.SignInForm(components.RenderAuthError(&[]utils.AuthError{{Err: err}})).Render(r.Context(), w)
			return
		}

		cred := utils.SignInCredentials{
			User: r.FormValue("user"),
			Pass: r.FormValue("password"),
		}

		user, errs := cred.SignIn()

		if errs != nil {
			components.SignInForm(components.RenderAuthError(errs)).Render(r.Context(), w)
			return
		}

		access, refresh, err := utils.NewTokenPair(
			&utils.JwtParams{
				Username: user.Username,
				UserId:   user.ID,
				Family:   uuid.New(),
			})

		if err != nil {
			components.SignInForm(components.RenderAuthError(errs)).Render(r.Context(), w)
			return
		}

		utils.SetTokenCookies(w, access, refresh)
		HandleRedirect(w, r, fmt.Sprintf("/profile/%s", claims.Username), http.StatusSeeOther, nil)
		return
	}

	if r.Method == "GET" {
		components.SignIn().Render(r.Context(), w)
		return
	}
}

func SignOut(w http.ResponseWriter, r *http.Request) {
	utils.DeleteJwtCookies(w)
	HandleRedirect(w, r, "/", http.StatusSeeOther, utils.ErrUserSignedOut)
}
