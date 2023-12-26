package routes

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"net/mail"

	"github.com/a-h/templ"
	"github.com/google/uuid"
	"github.com/nathan-hello/htmx-template/src/components"
	"github.com/nathan-hello/htmx-template/src/sqlc"
	"github.com/nathan-hello/htmx-template/src/utils"
	"github.com/nedpals/supabase-go"
)

type Credentials struct {
	username string
	password string
	email    string
}

func (c *Credentials) Validate() bool {
	user := len(c.username) > 3
	pass := len(c.password) > 4
	_, err := mail.ParseAddress(c.email)
	return user && pass && err == nil

}

func SignUp(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/signin" {
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

	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		cred := Credentials{
			username: r.FormValue("username"),
			password: r.FormValue("password"),
			email:    r.FormValue("email"),
		}

		if ok := cred.Validate(); !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		url := utils.Env().SUPABASE_POSTGREST_URL
		publicKey := utils.Env().SUPABASE_PUBLIC_KEY
		client := supabase.CreateClient(url, publicKey)

		user, err := client.Auth.SignUp(ctx, supabase.UserCredentials{
			Email:    cred.email,
			Password: cred.password,
		})

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var profileArgs sqlc.InsertProfileParams
		profileArgs.UserID, _ = uuid.Parse(user.ID)
		profileArgs.Username = cred.username

		existing, err := conn.SelectProfileByUsername(ctx, profileArgs.Username)
		if err != nil {
			if err != sql.ErrNoRows || existing.ID > 0 {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		prof, err := conn.InsertProfile(ctx, profileArgs)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/profile/%v", prof.ID), http.StatusSeeOther)
		return

	}

	if r.Method == "GET" {

		response, err := templ.ToGoHTML(ctx, components.SignIn())

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write([]byte(response))
		return

	}
}
