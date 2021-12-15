package lyrics

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/agnivade/levenshtein"
	"github.com/boratanrikulu/lirik.app/pkg/lyrics/constants"
	"github.com/gocolly/colly/v2"
	"golang.org/x/net/html"
)

type secondSource struct {
	Response struct {
		Hits []struct {
			Result struct {
				URL           string `json:"url"`
				Title         string `json:"title"`
				PrimaryArtist struct {
					Name string `json:"name"`
				} `json:"primary_artist"`
			} `json:"result"`
		} `json:"hits"`
	} `json:"response"`
}

func newSecondSource() *secondSource {
	return &secondSource{}
}

func (f *secondSource) GetLyrics(artistName string, songName string) (found bool, lyrics Lyrics) {
	u, _ := url.Parse(constants.SecondSourceAPI)
	q, _ := url.ParseQuery(u.RawQuery)

	q.Add("q", "\""+songName+" "+artistName+"\"")
	u.RawQuery = q.Encode()

	auth := "Bearer " + os.Getenv("GENIUS_ACCESS")
	req, _ := http.NewRequest("GET", fmt.Sprint(u), nil)
	req.Header.Set("Authorization", auth)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return false, lyrics
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&f)
	if err != nil {
		return false, lyrics
	}
	if len(f.Response.Hits) != 0 {
		gUrl := ""
		for _, value := range f.Response.Hits {
			s := strings.ToLower(songName)
			a := strings.ToLower(artistName)
			resultSong := strings.ToLower(value.Result.Title)
			resultArtist := strings.ToLower(value.Result.PrimaryArtist.Name)
			resultSong = songRegex(resultSong)

			// Second source use "’" for "'".
			s = strings.ReplaceAll(s, "'", "’")
			a = strings.ReplaceAll(a, "'", "’")

			if levenshtein.ComputeDistance(resultSong, s) <= 3 && strings.Contains(resultArtist, a) {
				gUrl = value.Result.URL
				break
			}
		}
		if gUrl == "" {
			return false, lyrics
		}

		c := colly.NewCollector()
		c.OnHTML(`#lyrics-root div[data-lyrics-container="true"]`, func(e *colly.HTMLElement) {
			lyrics.Lines = append(lyrics.Lines, lines(e.DOM)...)
		})

		c.Visit(gUrl)
	}

	if len(lyrics.Lines) == 0 {
		return false, lyrics
	}
	lyrics.Source = constants.SecondSourceBare
	return true, lyrics
}

func lines(s *goquery.Selection) []string {
	var resp []string

	var f func(*html.Node)
	var brFlag bool
	var tempLine string
	f = func(n *html.Node) {
		if n.Type == html.TextNode {
			tempLine += n.Data
			brFlag = false
		} else if n.Data == "br" {
			if brFlag {
				resp = append(resp, "")
				brFlag = false
			} else {
				brFlag = true
				resp = append(resp, tempLine)
				tempLine = ""
			}
		}
		if n.FirstChild != nil {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c)
			}
		}
	}
	for _, n := range s.Nodes {
		f(n)
		resp = append(resp, "")
	}

	return resp
}
