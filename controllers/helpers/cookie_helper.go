package helpers

import (
	"net/http"

	"github.com/boratanrikulu/s-lyrics/models"
)

func SetTokenCookies(w http.ResponseWriter, r models.RefreshAndAccessTokens) {
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

func SetStateCookie(w http.ResponseWriter, r models.Authorization) {
	cookie := http.Cookie {
		Name: "State",
		Value: r.Request.State,
	}
	http.SetCookie(w, &cookie)
}
