package routes

import (
	"net/http"
	"time"

	"github.com/nathan-hello/htmx-template/src/utils"
)

func SetTokenCookies(w http.ResponseWriter, a string, r string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    a,
		Expires:  time.Now().Add(utils.Env().REFRESH_EXPIRY_TIME), // Access token expiry time
		Secure:   true,
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    r,
		Expires:  time.Now().Add(utils.Env().REFRESH_EXPIRY_TIME), // Refresh token expiry time, e.g., 7 days
		Secure:   true,
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	})
}

func GetTokenFromCookie(r *http.Request) (string, string, error) {
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

func ValidateOrRefreshPairFromCookies(w http.ResponseWriter, r *http.Request) bool {

	access, err := r.Cookie("access_token")
	if err != nil {
		DeleteJwtCookies(w)
		return false
	}

	refresh, err := r.Cookie("refresh_token")
	if err != nil {
		DeleteJwtCookies(w)
		return false
	}

	vAccess, vRefresh, err := utils.ValidatePairOrRefresh(access.Value, refresh.Value)

	if err != nil {
		DeleteJwtCookies(w)
		return false
	}

	SetTokenCookies(w, vAccess, vRefresh)
	return true

}
