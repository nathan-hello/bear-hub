package routes

import (
	"log"
	"net/http"

	"github.com/nathan-hello/htmx-template/src/auth"
	"github.com/nathan-hello/htmx-template/src/components"
	"github.com/nathan-hello/htmx-template/src/utils"
)

func Auth(w http.ResponseWriter, r *http.Request) {
	state := utils.GetClientState(r)
	if state.IsAuthed {
		http.Redirect(w, r, "/profile/"+state.Username, http.StatusSeeOther)
		return
	}
	components.SignUp(state, auth.SignUp{}).Render(r.Context(), w)
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	state := utils.GetClientState(r)
	if state.IsAuthed {
		http.Redirect(w, r, "/profile/"+state.Username, http.StatusSeeOther)
		return
	}
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
		log.Println("START SIGNUP")
		components.SignUp(state, action).Render(r.Context(), w)
		log.Println("END SIGNUP")
		return
	}
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	state := utils.GetClientState(r)
	if state.IsAuthed {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	action := auth.SignIn{}

	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			action.MiscErrs = append(action.MiscErrs, err.Error())
			components.SignInForm(action).Render(r.Context(), w)
			return
		}

		action.UserOrEmail = r.FormValue("user")
		action.Password = r.FormValue("password")

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
