package routes

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/google/uuid"
	"github.com/nathan-hello/htmx-template/src/components"
	"github.com/nathan-hello/htmx-template/src/sqlc"
	"github.com/nathan-hello/htmx-template/src/utils"
	"github.com/nedpals/supabase-go"
)

func SignUp(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/signup" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	ctx := context.Background()
	db, err := sql.Open("postgres", utils.Env().DB_URI)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	conn := sqlc.New(db)

	returnFormWithErrors := func(errs *utils.SignUpErrors) {
		response, err := templ.ToGoHTML(ctx, components.SignInForm(*errs))

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write([]byte(response))
		return
	}

	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		cred := utils.Credentials{
			Username: r.FormValue("username"),
			Password: r.FormValue("password"),
			PassConf: r.FormValue("password-confirm"),
			Email:    r.FormValue("email"),
		}

		errs := cred.Validate()

		if len(errs.ErrsStr) > 0 {
			returnFormWithErrors(errs)
			return
		}

		errs = cred.ValidateEmailUsername()

		if len(errs.ErrsStr) > 0 {
			returnFormWithErrors(errs)
			return
		}

		url := utils.Env().SUPABASE_POSTGREST_URL
		publicKey := utils.Env().SUPABASE_PUBLIC_KEY
		client := supabase.CreateClient(url, publicKey)

		user, err := client.Auth.SignUp(ctx, supabase.UserCredentials{
			Email:    cred.Email,
			Password: cred.Password,
		})

		if err != nil {
			returnFormWithErrors(cred.CustomErrorMessage("Internal Server Error - 125632"))
			return
		}

		var profileArgs sqlc.InsertProfileParams
		profileArgs.UserID, _ = uuid.Parse(user.ID)
		profileArgs.Username = cred.Username
		prof, err := conn.InsertProfile(ctx, profileArgs)

		if err != nil {
			returnFormWithErrors(cred.CustomErrorMessage("Internal Server Error - 153294"))
			return
		}

		w.Header().Set("HX-Redirect", fmt.Sprintf("/profile/%v", prof.ID))
		return

	}

	if r.Method == "GET" {

		defaultFormErr := utils.FormErr{BorderColor: "border-blue-500"}
		form := utils.SignUpErrors{
			Email:    defaultFormErr,
			Username: defaultFormErr,
			Password: defaultFormErr,
			PassConf: defaultFormErr,
		}

		response, err := templ.ToGoHTML(ctx, components.SignIn(form))

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write([]byte(response))
		return

	}
}
