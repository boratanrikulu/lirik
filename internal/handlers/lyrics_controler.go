package handlers

import (
	"html/template"
	"net/http"

	"github.com/boratanrikulu/lirik.app/internal/handlers/helpers"
	"github.com/boratanrikulu/lirik.app/internal/models"
	"github.com/boratanrikulu/lirik.app/pkg/lyrics"
)

type LyricsPageData struct {
	IsAvaible bool
	Artist    models.Artist
	Song      models.Song
}

func LyricsGet(w http.ResponseWriter, r *http.Request) {
	artistName := r.URL.Query().Get("artistName")
	songName := r.URL.Query().Get("songName")
	albumImage := r.URL.Query().Get("albumImage")

	finder := lyrics.NewFinder()
	found, l := finder.GetLyrics(artistName, songName)
	pageData := LyricsPageData{
		IsAvaible: found,
		Artist: models.Artist{
			Name: artistName,
		},
		Song: models.Song{
			Name:       songName,
			Lyrics:     l,
			AlbumImage: albumImage,
		},
	}

	files := helpers.GetTemplateFiles("./views/songs.html")
	tmpl := template.Must(template.ParseFiles(files...))
	_ = tmpl.Execute(w, pageData)
}
