package models

import (
	"github.com/gocolly/colly/v2"
	"strings"
)

type Lyric struct {
	Lines []string
}

// Public Methods

func (l Lyric) GetLyric(artistName string, songName string) Lyric {
	c := colly.NewCollector()

	c.OnHTML("table#artistsonglist td.songName a[href]", func(e *colly.HTMLElement) {
		if strings.ToLower(strings.TrimSpace(e.Text)) == strings.ToLower(strings.TrimSpace(songName)) {
			link := "https://lyricstranslate.com/" + e.Attr("href")
			c.Visit(link)
		}
	})

	c.OnHTML(".ltf .par div", func(e *colly.HTMLElement) {
		l.Lines = append(l.Lines, e.Text)
	})

	path := "https://lyricstranslate.com/en/" + strings.Join(strings.Fields(artistName), "-") + "-lyrics.html"

	c.Visit(path)
	return l
}
