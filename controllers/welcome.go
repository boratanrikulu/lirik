package controllers

import (
	"html/template"
	"net/http"

	"github.com/boratanrikulu/s-lyrics/models"
)

// Page Datas

type LyricPageData struct {
	Artist models.Artist
	Song   models.Song
}

type WelcomePageData struct {
	SpotifyAuthLink string
}

// Public Methods

func WelcomeGet(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./views/welcome.html"))
	spotify := new(models.Spotify)

	_ = tmpl.Execute(w, WelcomePageData{
		SpotifyAuthLink: spotify.GetSpotifyAuthLink(),
	})
}

func LyricGet(w http.ResponseWriter, r *http.Request) {
	artistName := r.URL.Query().Get("artistName")
	songName := r.URL.Query().Get("songName")
	lyric := new(models.Lyric)

	pageData := LyricPageData{
		Artist: models.Artist{
			Name: artistName,
		},
		Song: models.Song{
			Name:  songName,
			Lyric: lyric.GetLyric(artistName, songName),
		},
	}

	tmpl := template.Must(template.ParseFiles("./views/songs.html"))
	_ = tmpl.Execute(w, pageData)
}
