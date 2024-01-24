package routes

import (
	"net/http"
)

func Redirect500(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/?500=true", http.StatusSeeOther)
}

func Redirect404(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/?404=true", http.StatusSeeOther)
}

func Redirect401(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/?unauthorized=true", http.StatusSeeOther)
}

func RedirectSignOut(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/?signedout=true", http.StatusSeeOther)
}

func RedirectToSignIn(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/signin", http.StatusSeeOther)
}

func Alert(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/alert/signout" {
		w.Write([]byte("You've been signed out"))
	}

	if r.URL.Path == "/alert/unauthorized" {
		w.Write([]byte("You're not logged in"))
	}

	if r.URL.Path == "/alert/404" {
		w.Write([]byte("404 Not Found"))
	}

	if r.URL.Path == "/alert/500" {
		w.Write([]byte("500 Internal Server Error"))
	}

}
