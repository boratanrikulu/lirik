package helpers

import (
	"net/http"
	"time"

	"github.com/boratanrikulu/s-lyrics/models"
)

func SetTokenCookies(r models.RefreshAndAccessTokens, w http.ResponseWriter) {
	cookies := []http.Cookie {
		http.Cookie {
			Name: "AccessToken",
			Value: r.Response.AccessToken,
		},
		http.Cookie {
			Name: "RefreshToken",
			Value: r.Response.RefreshToken,
		},
	}
	for _, cookie := range cookies {
		http.SetCookie(w, &cookie)
	}
}

func SetStateCookie(r models.Authorization, w http.ResponseWriter) {
	cookie := http.Cookie {
		Name: "State",
		Value: r.Request.State,
	}
	http.SetCookie(w, &cookie)
}

func ClearCookies(w http.ResponseWriter, r *http.Request) {
	// Clears all cookies.
	for _, cookie := range r.Cookies() {
		cookie.Value = ""
		cookie.Expires = time.Unix(0, 0)
		http.SetCookie(w, cookie)
	}
}
