package controllers

import (
	"fmt"
	"github.com/boratanrikulu/s-lyrics/models"
	"html/template"
	"net/http"
)

// Page Datas

type LyricPageData struct {
	Lyric models.Lyric
}

// Public Methods

func LyricsGet(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Lyric controller")
	lyric := new(models.Lyric)
	lyric.SongName = r.URL.Query().Get("songName")
	lyric.ArtistName = r.URL.Query().Get("artistName")
	err := lyric.GetLyric()
	if err != nil {
		panic("ok")
	}

	tmpl := template.Must(template.ParseFiles("./views/lyrics.html"))
	data := LyricPageData{Lyric: *lyric}
	_ = tmpl.Execute(w, data)
}