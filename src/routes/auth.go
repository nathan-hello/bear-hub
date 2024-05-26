package routes

import (
	"context"
	"net/http"

	"github.com/nathan-hello/htmx-template/src/auth"
	"github.com/nathan-hello/htmx-template/src/components"
	"github.com/nathan-hello/htmx-template/src/utils"
)

func sendArbitraryError(err error, ctx context.Context, w http.ResponseWriter) {
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.ClaimsContextKey).(auth.CustomClaims)
	if ok {
		http.Redirect(w, r, "/profile/"+claims.Username, http.StatusSeeOther)
		return
	}
	state := components.AuthSignUpState{
		ClientState: components.ClientState{IsAuthed: ok},
	}

	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			sendArbitraryError(utils.ErrParseForm, r.Context(), w)
		}

		state.Username = r.FormValue("username")
		state.Password = r.FormValue("password")
		state.PassConf = r.FormValue("password-confirmation")
		state.Email = r.FormValue("email")

		var signInAction = auth.AuthSignUp{
			State: &state,
		}

                user := signInAction.SignUp()

		access, refresh, err := auth.NewTokenPair(
			&auth.JwtParams{
				Username: user.Username,
				UserId:   user.ID,
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
		http.Redirect(w, r, "/profile/%s"+claims.Username, http.StatusSeeOther)
		return
	}

	if r.Method == "GET" {
		components.SignIn(state).Render(r.Context(), w)
		return
	}
}

func SignOut(w http.ResponseWriter, r *http.Request) {
	auth.DeleteJwtCookies(w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
