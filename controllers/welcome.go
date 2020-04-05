package controllers

import (
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
	tmpl := template.Must(template.ParseFiles("./views/welcome.html"))
	spotify := new(models.Spotify)
	spotify.InitSecrets()

	_ = tmpl.Execute(w, WelcomePageData{
		SpotifyAuthLink: spotify.GetRequestAuthorizationLink(),
	})
}
