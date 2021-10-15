package lyrics

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/boratanrikulu/lirik.app/pkg/lyrics/constants"
	"github.com/gocolly/colly/v2"
)

type firstSource struct{}

func newFirstSource() *firstSource {
	return &firstSource{}
}

func (f *firstSource) GetLyrics(artistName string, songName string) (found bool, lyrics Lyrics) {
	c := colly.NewCollector()

	// Search and find the lyric page url.
	a := url.PathEscape("\"" + artistName + "\"")
	s := url.PathEscape("\"" + songName + "\"")
	url := fmt.Sprintf("%s/en/songs/0/%s/%s", constants.FirstSource, a, s)
	url = strings.ReplaceAll(url, "%", "%25")
	songUrl := ""
	c.OnHTML(".ltsearch-results-line tbody tr", func(e *colly.HTMLElement) {
		e.ForEach("td:nth-child(2)", func(_ int, e *colly.HTMLElement) {
			s = strings.ToLower(songName)
			resultSong := strings.TrimSpace(strings.ToLower(e.Text))
			resultSong = songRegex(resultSong)

			if resultSong == s {
				songUrl = e.ChildAttr("a[href]", "href")
			}
		})
	})
	c.Visit(url)
	if songUrl == "" {
		return false, lyrics
	}

	// Go to song page and take lyrics.
	cc := colly.NewCollector()
	cc.OnHTML("div#song-body .ltf .par div, .emptyline", func(e *colly.HTMLElement) {
		line := strings.Trim(strings.TrimSpace(e.Text), "\n")
		lyrics.Lines = append(lyrics.Lines, line)
	})
	cc.OnHTML(".langsmall-song span.langsmall-languages", func(e *colly.HTMLElement) {
		if strings.TrimSpace(e.Text) != "" {
			lyrics.Language = strings.Trim(strings.TrimSpace(e.Text), "\n")
		}
	})
	songUrl = constants.FirstSource + songUrl
	cc.Visit(songUrl)

	if len(lyrics.Lines) == 0 {
		return false, Lyrics{}
	}

	lyrics.Source = constants.FirstSourceBare
	lyrics.Translates = f.getTranslations(songUrl)
	return true, lyrics
}

func (f *firstSource) getTranslations(url string) (translates []Translate) {
	c := colly.NewCollector()

	var addedTranslations []string
	c.OnHTML("div.song-node div.masonry-grid div.song-list-translations-list a[href]", func(e *colly.HTMLElement) {
		cc := colly.NewCollector()
		cc.OnHTML("div.translate-node-text", func(e *colly.HTMLElement) {
			language := e.ChildText("div.langsmall-song span.mobile-only-inline")
			if !contains(addedTranslations, language) {
				translate := Translate{}
				translate.Language = language
				e.ForEach("div.authorsubmitted a[href]", func(c int, e *colly.HTMLElement) {
					if c == 0 {
						translate.Author.Name = e.Text
						translate.Author.Href = e.Attr("href")
					}
				})
				if translate.Language != "" {
					translate.Title = e.ChildText("h2.title-h2")
					e.ForEach(".ltf .par div, .emptyline", func(_ int, e *colly.HTMLElement) {
						line := strings.Trim(strings.TrimSpace(e.Text), "\n")
						translate.Lines = append(translate.Lines, line)
					})
					translates = append(translates, translate)
					addedTranslations = append(addedTranslations, translate.Language)
				}
			}
		})

		if contains(constants.AllowedTranslationLanguages, e.Text) {
			cc.Visit(fmt.Sprintf("%s/%s", constants.FirstSource, e.Attr("href")))
		}
	})

	c.Visit(url)
	return translates
}
