package models

import (
	"github.com/gocolly/colly/v2"
	"regexp"
	"net/url"
	"strings"
	"fmt"
)

type Lyric struct {
	Lines []string
	IsAvaible bool
}

// Public Methods

func (l Lyric) GetLyric(artistName string, songName string) Lyric {
	// Removes values after "(..." or "-...". from song name.
	re := regexp.MustCompile(`[-(].+`)
	songName = re.ReplaceAllString(songName, "")

	c := colly.NewCollector()

	artistName = url.PathEscape("\"" + artistName + "\"")
	songName = url.PathEscape("\"" + songName + "\"")
	url := "https://lyricstranslate.com/en/songs/0/" + artistName + "/" + songName
	// TODO fix this issue
	url = strings.ReplaceAll(url, "%", "%25")
	counter := 0
	c.OnHTML(".ltsearch-results-line tbody tr td a[href]", func(e *colly.HTMLElement) {
		counter++
		if counter == 2 {
			// That means it is song value.
			c.Visit("https://lyricstranslate.com/" + e.Attr("href"))
		}
	})

	c.OnHTML("div#song-body .ltf .par div, .emptyline", func(e *colly.HTMLElement) {
		l.Lines = append(l.Lines, e.Text)
	})

	c.Visit(fmt.Sprint(url))

	if len(l.Lines) != 0 {
		l.IsAvaible = true
	}
	return l
}
