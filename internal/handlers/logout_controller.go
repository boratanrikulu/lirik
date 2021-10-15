package handlers

import (
	"net/http"

	"github.com/boratanrikulu/lirik.app/internal/handlers/helpers"
)

func LogoutGet(w http.ResponseWriter, r *http.Request) {
	helpers.ClearCookies(w, r)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
