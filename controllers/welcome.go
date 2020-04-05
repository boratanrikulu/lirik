package controllers

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/boratanrikulu/s-lyrics/models"
	"github.com/gocolly/colly/v2"
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

	pageData := LyricPageData{
		Artist: models.Artist{
			Name: artistName,
		},
		Song: models.Song{
			Name:  songName,
			Lyric: getLyric(artistName, songName),
		},
	}

	tmpl := template.Must(template.ParseFiles("./views/songs.html"))
	_ = tmpl.Execute(w, pageData)
}

// Private Methods

func getLyric(artistName string, songName string) models.Lyric {
	c := colly.NewCollector()

	lyric := models.Lyric{}

	c.OnHTML("table#artistsonglist td.songName a[href]", func(e *colly.HTMLElement) {
		if strings.ToLower(strings.TrimSpace(e.Text)) == strings.ToLower(strings.TrimSpace(songName)) {
			link := "https://lyricstranslate.com/" + e.Attr("href")
			c.Visit(link)
		}
	})

	c.OnHTML(".ltf .par div", func(e *colly.HTMLElement) {
		lyric.Lines = append(lyric.Lines, e.Text)
	})

	path := "https://lyricstranslate.com/en/" + strings.Join(strings.Fields(artistName), "-") + "-lyrics.html"

	c.Visit(path)
	return lyric
}
