package controllers

import (
	"github.com/boratanrikulu/s-lyrics/models"
	"html/template"
	"net/http"
)

// Page Datas

type LyricPageData struct {
	Artist models.Artist
	Song   models.Song
}

// Public Methods

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
