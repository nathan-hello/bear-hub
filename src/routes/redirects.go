package routes

import (
	"fmt"
	"net/http"
)

func Redirect500(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/?500=/", http.StatusSeeOther)
}

func Redirect404(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/?404=/", http.StatusSeeOther)
}

func Redirect401(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/?unauthorized=/", http.StatusSeeOther)
}

func RedirectSignOut(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/?signedout=/", http.StatusSeeOther)
}

func RedirectToSignIn(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/signin", http.StatusSeeOther)
}

func RedirectToProfile(w http.ResponseWriter, r *http.Request, username string) {
	http.Redirect(w, r, fmt.Sprintf("/profile/%v", username), http.StatusSeeOther)
}
