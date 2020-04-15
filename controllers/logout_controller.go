package controllers

import (
	"github.com/boratanrikulu/s-lyrics/controllers/helpers"
	"net/http"
)

// Public Methods

func LogoutGet(w http.ResponseWriter, r *http.Request) {
	helpers.ClearCookies(w, r)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
