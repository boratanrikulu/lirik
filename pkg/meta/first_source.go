package meta

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/boratanrikulu/lirik.app/pkg/meta/constants"
	"github.com/gocolly/colly/v2"
)

type firstSource struct{}

func newFirstSource() *firstSource {
	return &firstSource{}
}

func (f *firstSource) GetMeta(artistName string, albumName string) (found bool, meta Meta) {
	albumName = albumRegex(albumName)
	u, _ := url.Parse(constants.FirstSource + "/search")
	q, _ := url.ParseQuery(u.RawQuery)

	q.Add("title", albumName)
	q.Add("artist", artistName)
	q.Add("type", "all")
	u.RawQuery = q.Encode()

	albumLink := ""
	c := colly.NewCollector()
	c.OnHTML(
		`#search_results div[data-object-type]:first-child a[href].search_result_title:first-child`,
		func(e *colly.HTMLElement) {
			albumLink = e.Attr("href")
		},
	)
	c.Visit(fmt.Sprint(u))
	if albumLink == "" {
		return false, meta
	}

	cc := colly.NewCollector()
	cc.OnHTML(
		`div.profile`,
		func(e *colly.HTMLElement) {
			texts := e.ChildTexts("div")
			for i, text := range texts {
				if text == "Genre:" && i+1 < len(texts) {
					meta.Genre = texts[i+1]
				}
				if text == "Style:" && i+1 < len(texts) {
					meta.Style = texts[i+1]
				}
			}
		},
	)
	cc.OnHTML(
		`table.table_1fWaB tbody tr`,
		func(e *colly.HTMLElement) {
			if e.ChildText("th") == "Genre:" {
				meta.Genre = e.ChildText("td")
			}
			if e.ChildText("th") == "Style:" {
				meta.Style = e.ChildText("td")
			}
		},
	)
	cc.Visit(constants.FirstSource + albumLink)
	if meta.Genre == "" {
		return false, meta
	}

	return true, meta
}

func albumRegex(album string) string {
	regexList := []string{
		` - .+`,
		`(?i)\(.*\)`,
	}

	for _, value := range regexList {
		re := regexp.MustCompile(value)
		album = re.ReplaceAllString(album, "")
	}

	album = strings.TrimSpace(album)
	return album
}
