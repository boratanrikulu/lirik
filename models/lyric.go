package models

import (
	"github.com/gocolly/colly/v2"
	"regexp"
	"net/url"
	"strings"
	"fmt"
)

type Lyric struct {
	Lines 		[]string
	IsAvaible	bool
	Translates 	[]Translate
}

type Translate struct {
	Language		string
	Title			string
	Lines			[]string
}

// Public Methods

func (l Lyric) GetLyric(artistName string, songName string) Lyric {
	// Removes values after "(..." or "-...". from song name.
	re := regexp.MustCompile(`[-(].+`)
	songName = re.ReplaceAllString(songName, "")

	c := colly.NewCollector()

	// Search lyric for the song.
	artistName = url.PathEscape("\"" + artistName + "\"")
	songName = url.PathEscape("\"" + songName + "\"")
	url := "https://lyricstranslate.com/en/songs/0/" + artistName + "/" + songName
	// TODO fix this issue
	url = strings.ReplaceAll(url, "%", "%25")
	counter := 0
	songUrl := ""
	c.OnHTML(".ltsearch-results-line tbody tr td a[href]", func(e *colly.HTMLElement) {
		counter++
		if counter == 2 {
			// That means it is song value.
			songUrl = "https://lyricstranslate.com/" + e.Attr("href")
			c.Visit(songUrl)
		}
	})

	// Song lyric page.
	c.OnHTML("div#song-body .ltf .par div, .emptyline", func(e *colly.HTMLElement) {
		l.Lines = append(l.Lines, e.Text)
	})

	c.Visit(fmt.Sprint(url))

	if len(l.Lines) != 0 {
		l.IsAvaible = true
		// Gets avaible translations for the song.
		getTranslations(&l, songUrl)
	}
	return l
}

// Private Methods

func getTranslations(l *Lyric, url string) {
	fmt.Println(url)
	c := colly.NewCollector()

	// Translation list for the song.
	c.OnHTML("div.song-node-info li.song-node-info-translate a[href]", func(e *colly.HTMLElement) {
		// TODO
		// Fix more-then-one translate issue.
		c.Visit("https://lyricstranslate.com/" + e.Attr("href"))
	})

	// Lyric translations for the song.
	c.OnHTML("div.translate-node-text", func(e *colly.HTMLElement) {
		translate := Translate{}
		translate.Language = e.ChildText("div.langsmall-song span.mobile-only-inline")
		if translate.Language != "" {
			translate.Title = e.ChildText("h2.title-h2")
			e.ForEach(".ltf .par div, .emptyline", func(_ int, e *colly.HTMLElement) {
				translate.Lines = append(translate.Lines, e.Text)
			})
			l.Translates = append(l.Translates, translate)
		}
	})

	c.Visit(url)
}
