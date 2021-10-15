package helpers

import (
	"net/http"
	"time"

	"github.com/boratanrikulu/lirik.app/pkg/spotify"
)

func UpdateTokenCookies(s spotify.SpotifyAuthorizatied, w http.ResponseWriter) {
	cookies := []http.Cookie{
		{
			Name:     "AccessToken",
			Value:    s.Authorization.AccessToken,
			MaxAge:   s.Authorization.ExpiresIn,
			SameSite: http.SameSiteLaxMode,
			HttpOnly: false,
		},
		{
			Name:     "Me",
			Value:    s.UserMe,
			MaxAge:   s.Authorization.ExpiresIn,
			SameSite: http.SameSiteLaxMode,
			HttpOnly: false,
		},
	}

	if s.Authorization.RefreshToken != "" {
		oneMonth := time.Hour * 24 * 30
		cookie := http.Cookie{
			Name:     "RefreshToken",
			Value:    s.Authorization.RefreshToken,
			MaxAge:   int(oneMonth.Seconds()),
			SameSite: http.SameSiteLaxMode,
			HttpOnly: false,
		}
		cookies = append(cookies, cookie)
	}

	for _, cookie := range cookies {
		http.SetCookie(w, &cookie)
	}
}

func SetStateCookie(state string, w http.ResponseWriter) {
	cookie := http.Cookie{
		Name:     "State",
		Value:    state,
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
}

func ClearCookies(w http.ResponseWriter, r *http.Request) {
	for _, cookie := range r.Cookies() {
		cookie.Value = ""
		cookie.Expires = time.Unix(0, 0)
		http.SetCookie(w, cookie)
	}
}
