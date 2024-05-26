package routes

import (
	"net/http"

	"github.com/nathan-hello/htmx-template/src/auth"
	"github.com/nathan-hello/htmx-template/src/components"
)

func Auth(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.ClaimsContextKey).(auth.CustomClaims)
	if ok {
		http.Redirect(w, r, "/profile/"+claims.Username, http.StatusSeeOther)
		return
	}
	state := components.ClientState{IsAuthed: ok}
       components.SignUp(state, auth.SignUp{}).Render(r.Context(), w)
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.ClaimsContextKey).(auth.CustomClaims)
	if ok {
		http.Redirect(w, r, "/profile/"+claims.Username, http.StatusSeeOther)
		return
	}
	state := components.ClientState{IsAuthed: ok}
        action := auth.SignUp{}

	if r.Method == "POST" {
		err := r.ParseForm()

		if err != nil {
                        action.MiscErrs = append(action.MiscErrs, err.Error())
                        components.SignUpForm(action).Render(r.Context(), w)
                        return
		}

                action.Username = r.FormValue("username")
		action.Password = r.FormValue("password")
		action.PassConf = r.FormValue("password-confirmation")
		action.Email = r.FormValue("email")
                user := action.SignUp()

                if user == nil {
                        components.SignUpForm(action).Render(r.Context(), w) // There will be errors in action.RenderErrs() if user == nil
                        return
                }

		access, refresh, err := auth.NewTokenPair(
			&auth.JwtParams{
				Username: user.Username,
				UserId:   user.ID,
			})

                if err != nil {
                        action.MiscErrs = append(action.MiscErrs, err.Error())
                        components.SignUpForm(action).Render(r.Context(), w)
                        return
                }

		auth.SetTokenCookies(w, access, refresh)
		w.Header().Set("HX-Redirect", "/")
		return

	}

	if r.Method == "GET" {
		components.SignUp(state, action).Render(r.Context(), w)
		return
	}
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value(auth.ClaimsContextKey).(*auth.CustomClaims)
	if ok {
		w.Header().Set("HX-Redirect", "/")
		return
	}
	state := components.ClientState{ IsAuthed: ok, }
        action := auth.SignIn{}

	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
                        action.MiscErrs = append(action.MiscErrs, err.Error())
                        components.SignInForm(action).Render(r.Context(), w)
                        return
		}

			action.UserOrEmail= r.FormValue("user")
			action.Password= r.FormValue("password")

                user := action.SignIn()

                if user == nil {
                        components.SignInForm(action).Render(r.Context(), w)
                        return
                }

		access, refresh, err := auth.NewTokenPair(
			&auth.JwtParams{
				Username: user.Username,
				UserId:   user.ID,
			})

                if err != nil {
                        action.MiscErrs = append(action.MiscErrs, err.Error())
                        components.SignInForm(action).Render(r.Context(), w)
                        return
                }


		auth.SetTokenCookies(w, access, refresh)
		http.Redirect(w, r, "/profile/%s"+user.Username, http.StatusSeeOther)
		return
	}

	if r.Method == "GET" {
		components.SignIn(state, action).Render(r.Context(), w)
		return
	}
}

func SignOut(w http.ResponseWriter, r *http.Request) {
	auth.DeleteJwtCookies(w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
