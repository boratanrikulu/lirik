package models

import (
	"github.com/gocolly/colly/v2"
	"strings"
	"fmt"
	"regexp"
)

type Lyric struct {
	Lines []string
}

// Public Methods

func (l Lyric) GetLyric(artistName string, songName string) Lyric {
	// Removes values after "(..." or "-...".
	re := regexp.MustCompile(`[-(].+`)
	songName = re.ReplaceAllString(songName, "")

	c := colly.NewCollector()

	c.OnHTML("table#artistsonglist td.songName a[href]", func(e *colly.HTMLElement) {
		eTextTrim := strings.ToLower(strings.TrimSpace(e.Text))
		songNameTrim := strings.ToLower(strings.TrimSpace(songName))
		if eTextTrim == songNameTrim {
			link := "https://lyricstranslate.com/" + e.Attr("href")
			c.Visit(link)
		}
	})

	c.OnHTML("div#song-body .ltf .par div, .emptyline", func(e *colly.HTMLElement) {
		l.Lines = append(l.Lines, e.Text)
	})

	url := "https://lyricstranslate.com/en/" +
				strings.Join(strings.Fields(artistName), "-") +
				"-lyrics.html"
	c.Visit(fmt.Sprint(url))
	return l
}
