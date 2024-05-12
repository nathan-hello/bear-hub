package auth

import (
	"net/http"
	"time"

	"github.com/nathan-hello/htmx-template/src/utils"
)

func SetTokenCookies(w http.ResponseWriter, a string, r string) {
        secure := true
        if utils.Env().MODE == "dev" {
                secure = false
        }
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    a,
		Expires:  time.Now().Add(utils.Env().REFRESH_EXPIRY_TIME),
		Secure:   secure,
		HttpOnly: secure,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,

	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    r,
		Expires:  time.Now().Add(utils.Env().REFRESH_EXPIRY_TIME),
		Secure:   secure,
		HttpOnly: secure,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	})
}

func GetJwtsFromCookie(r *http.Request) (string, string, error) {
	access, err := r.Cookie("access_token")
	if err != nil {
		return "", "", err
	}

	refresh, err := r.Cookie("refresh_token")
	if err != nil {
		return "", "", err
	}

	return access.Value, refresh.Value, nil
}

func DeleteCookie(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
}

func DeleteJwtCookies(w http.ResponseWriter) {
	DeleteCookie(w, "access_token")
	DeleteCookie(w, "refresh_token")
}

func ValidateJwtOrDelete(w http.ResponseWriter, r *http.Request) (string, bool) {
	access, err := r.Cookie("access_token")
	if err != nil {
                if err == http.ErrNoCookie {
                        return "", false
                }
                utils.PrintlnOnDevMode("delete access because r.Cookie err", err) 
		DeleteJwtCookies(w)
		return "", false
	}

	refresh, err := r.Cookie("refresh_token")
	if err != nil {
                if err == http.ErrNoCookie {
                        return "", false
                }
                utils.PrintlnOnDevMode("delete refresh because r.Cookie err", err) 
		DeleteJwtCookies(w)
		return "", false
	}

	vAccess, vRefresh, err := ValidatePairOrRefresh(access.Value, refresh.Value)

	if err != nil {
                utils.PrintlnOnDevMode("delete vaccess/vrefresh", err) 
		DeleteJwtCookies(w)
		return "", false
	}

	SetTokenCookies(w, vAccess, vRefresh)
	return vAccess, true
}
