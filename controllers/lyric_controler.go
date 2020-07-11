package controllers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/boratanrikulu/s-lyrics/controllers/helpers"
	"github.com/boratanrikulu/s-lyrics/models"
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

	if !l.IsAvaible {
		log.Printf("[NOT FOUND] \"%v by %v\"", songName, artistName)
	} else {
		log.Printf("[FOUND] \"%v by %v\"", songName, artistName)
	}

	files := helpers.GetTemplateFiles("./views/songs.html")
	tmpl := template.Must(template.ParseFiles(files...))
	_ = tmpl.Execute(w, pageData)
}
