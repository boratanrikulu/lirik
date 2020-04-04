package controllers

import (
  "strings"
  "net/url"
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
  spotifyAuthLink := spotifyAuthLink()
  tmpl, _ := template.ParseFiles("./views/welcome.html")
  tmpl.Execute(w, spotifyAuthLink)
}

func WelcomePost(w http.ResponseWriter, r *http.Request) {
  artistName := r.URL.Query().Get("artistName")
  songName := r.URL.Query().Get("songName")

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

func spotifyAuthLink() string {
  baseUrl, _ := url.Parse("https://accounts.spotify.com/")
  baseUrl.Path += "authorize"

  params := url.Values{}
  params.Add("client_id", "6f524a004e874120b42251c6c6d0e699")
  params.Add("response_type", "code")
  params.Add("redirect_uri", "http://localhost:3000/spotify")
  params.Add("scope", "user-read-currently-playing streaming user-read-playback-state")
  params.Add("state", "alskdjalskfnalsdkmalskdm")
  baseUrl.RawQuery = params.Encode()

  return baseUrl.String()
}

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
