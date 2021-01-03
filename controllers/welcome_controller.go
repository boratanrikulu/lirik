package controllers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/boratanrikulu/lirik.app/controllers/helpers"
	"github.com/boratanrikulu/lirik.app/models"
)

// Page Datas

type WelcomePageData struct {
	SpotifyAuthLink string
}

// Public Methods

func WelcomeGet(w http.ResponseWriter, r *http.Request) {
	// If there is an access or refresh token,
	// redirect to "/spotify" page.
	accessTokenCookie, _ := r.Cookie("AccessToken")
	refreshTokenCookie, _ := r.Cookie("RefreshToken")
	if accessTokenCookie != nil || refreshTokenCookie != nil {
		http.Redirect(w, r, "/spotify", http.StatusSeeOther)
		return
	}

	files := helpers.GetTemplateFiles("./views/welcome.html")
	tmpl := template.Must(template.ParseFiles(files...))
	spotify := new(models.Spotify)
	spotify.InitSecrets()

	authLink, err := spotify.GetRequestAuthorizationLink()
	helpers.SetStateCookie(spotify.Authorization, w)
	if err != nil {
		log.Println("Somethings are wrong. We are working on it.")
		helpers.ErrorPage([]string{"Somethings are wrong. We are working on it."}, w)
		return
	}
	_ = tmpl.Execute(w, WelcomePageData{
		SpotifyAuthLink: authLink,
	})
}
