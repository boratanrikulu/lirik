package controllers

import (
	"html/template"
	"net/http"

	"github.com/boratanrikulu/lirik.app/controllers/helpers"
	"github.com/boratanrikulu/lirik.app/models"
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
	albumImage := r.URL.Query().Get("albumImage")

	l := new(models.Lyric)
	l.GetLyricByCheckingDatabase(artistName, songName)
	count := 1
	for count < 3 && !l.IsAvaible {
		l.GetLyricByCheckingDatabase(artistName, songName)
		count++
	}

	pageData := LyricPageData{
		Artist: models.Artist{
			Name: artistName,
		},
		Song: models.Song{
			Name:       songName,
			Lyric:      *l,
			AlbumImage: albumImage,
		},
	}

	files := helpers.GetTemplateFiles("./views/songs.html")
	tmpl := template.Must(template.ParseFiles(files...))
	_ = tmpl.Execute(w, pageData)
}
