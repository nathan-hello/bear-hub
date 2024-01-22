package routes

import (
	"net/http"
	"time"

	"github.com/nathan-hello/htmx-template/src/utils"
)

func SetTokenCookies(w http.ResponseWriter, accessToken string, refreshToken string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Expires:  time.Now().Add(utils.Env().REFRESH_EXPIRY_TIME), // Access token expiry time
		Secure:   true,
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(utils.Env().REFRESH_EXPIRY_TIME), // Refresh token expiry time, e.g., 7 days
		Secure:   true,
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	})
}
