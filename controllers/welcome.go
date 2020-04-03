package controllers

import (
  "strings"
  "net/http"
  "html/template"

  "github.com/boratanrikulu/s-lyrics/models"
  "github.com/gocolly/colly/v2"
)

// Page Datas

type PageData struct {
  Artist models.Artist
  Song models.Song
}

// Public Methods

func WelcomeGet(w http.ResponseWriter, r *http.Request) {
  tmpl, _ := template.ParseFiles("./views/welcome.html")
  tmpl.Execute(w, nil)
}

func WelcomePost(w http.ResponseWriter, r *http.Request) {
  artistName := r.PostFormValue("artistName")
  songName := r.PostFormValue("songName")

  pageData := PageData{
    Artist: models.Artist{
      Name: artistName,
    },
    Song: models.Song{
      Name: songName,
      Lyric: getLyric(artistName, songName),
    },
  }

  tmpl, _ := template.ParseFiles("./views/songs.html")
  tmpl.Execute(w, pageData)
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
