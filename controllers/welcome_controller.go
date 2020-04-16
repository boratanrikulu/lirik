package controllers

import (
	"github.com/boratanrikulu/s-lyrics/controllers/helpers"
	"github.com/boratanrikulu/s-lyrics/models"
	"html/template"
	"net/http"
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
		panic("Something is wrong")
	}
	err = tmpl.Execute(w, WelcomePageData{
		SpotifyAuthLink: authLink,
	})
	if err != nil {
		panic("Something is wrong")
	}
}
