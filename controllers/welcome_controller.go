package controllers

import (
	"github.com/boratanrikulu/s-lyrics/models"
	"github.com/boratanrikulu/s-lyrics/controllers/helpers"
	"html/template"
	"net/http"
)

// Page Datas

type WelcomePageData struct {
	SpotifyAuthLink string
}

// Public Methods

func WelcomeGet(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./views/welcome.html"))
	spotify := new(models.Spotify)
	spotify.InitSecrets()

	authLink, err := spotify.GetRequestAuthorizationLink()
	helpers.SetStateCookie(w, spotify.Authorization)
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
