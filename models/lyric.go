package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/gocolly/colly/v2"
)

type Lyric struct {
	Lines      []string
	IsAvaible  bool
	Language   string
	Translates []Translate
	Source     string
}

type Translate struct {
	Language string
	Title    string
	Author   struct {
		Name string
		Href string
	}
	Lines []string
}

type Genius struct {
	Response struct {
		Hits []struct {
			Type   string `json:type`
			Result struct {
				URL string `json:"url"`
				Title string `json:"title"`
				PrimaryArtist struct {
					Name string `json:"name"`
				} `json:"primary_artist"`
			} `json:"result"`
		} `json:"hits"`
	} `json:"response"`
}

// Public Methods

func (l Lyric) GetLyric(artistName string, songName string) Lyric {
	// Removes values after " - ...". from song name.
	rere := regexp.MustCompile(` - .+`)
	// Remove all (...)s from song name.
	re := regexp.MustCompile(`\(.*?\)`)
	songName = rere.ReplaceAllString(songName, "")
	songName = re.ReplaceAllString(songName, "")
	// Trim spaces
	songName = strings.TrimSpace(songName)

	// Get from lyricstranslates.com
	getFromFirstSource(&l, artistName, songName)

	if !l.IsAvaible {
		// If there is no lyric on the first source,
		// then get it from genius.com
		getFromSecondSource(&l, artistName, songName)
		
	}

	// TODO remove "return" and user pointers.
	return l
}

// Private Methods

func getFromSecondSource(l *Lyric, artistName string, songName string) {
	u, _ := url.Parse("https://api.genius.com/search")
	q, _ := url.ParseQuery(u.RawQuery)

	q.Add("q", "\""+songName+" "+artistName+"\"")
	u.RawQuery = q.Encode()

	auth := "Bearer " + os.Getenv("GENIUS_ACCESS")
	req, _ := http.NewRequest("GET", fmt.Sprint(u), nil)
	req.Header.Set("Authorization", auth)

	// Sends the request.
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	genius := new(Genius)
	json.Unmarshal(body, &genius)

	if len(genius.Response.Hits) != 0 {
		geniusURL := ""
		for _, value := range genius.Response.Hits {
			s := strings.ToLower(songName)
			a := strings.ToLower(artistName)
			resultSong := strings.ToLower(value.Result.Title)
			resultArtist := strings.ToLower(value.Result.PrimaryArtist.Name)

			// Yes. But true.
			// Genius use "’" for "'".
			// Btw, How's the Heart?
			s = strings.ReplaceAll(s, "'", "’")
			a = strings.ReplaceAll(a, "'", "’")

			if resultSong == s && resultArtist == a {
				geniusURL = value.Result.URL
				break
			}
		}

		if geniusURL == "" {
			return
		}

		c := colly.NewCollector()

		c.OnHTML("div.song_body-lyrics div.lyrics p", func(e *colly.HTMLElement) {
			lines := strings.SplitAfter(e.Text, "\n")
			for _, line := range lines {
				if line == "\n" {
					line = ""
				}
				l.Lines = append(l.Lines, line)
			}
		})

		c.Visit(geniusURL)
	}

	if len(l.Lines) != 0 {
		l.Source = "genius.com"
		l.IsAvaible = true
	}
}

func getFromFirstSource(l *Lyric, artistName string, songName string) {
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
		cc := colly.NewCollector()

		// Song lyric page.
		cc.OnHTML("div#song-body .ltf .par div, .emptyline", func(e *colly.HTMLElement) {
			l.Lines = append(l.Lines, e.Text)
		})

		// Song's language.
		cc.OnHTML(".langsmall-song span.langsmall-languages", func(e *colly.HTMLElement) {
			if strings.TrimSpace(e.Text) != "" {
				l.Language = e.Text
			}
		})

		counter++
		if counter == 2 {
			// That means it is song value.
			songUrl = "https://lyricstranslate.com/" + e.Attr("href")
			// Visit song page.
			cc.Visit(songUrl)
		}
	})

	// Vist search page.
	c.Visit(fmt.Sprint(url))

	if len(l.Lines) != 0 {
		l.Source = "lyricstranslate.com"
		l.IsAvaible = true
		// Gets avaible translations for the song.
		getTranslations(l, songUrl)
	}
}

func getTranslations(l *Lyric, url string) {
	c := colly.NewCollector()

	allowedTranslationLanguages := "Turkish English Italian Swedish German French"
	// Translation list for the song.
	c.OnHTML("div.song-node-info li.song-node-info-translate a[href]", func(e *colly.HTMLElement) {
		cc := colly.NewCollector()

		// Lyric translations for the song.
		cc.OnHTML("div.translate-node-text", func(e *colly.HTMLElement) {
			translate := Translate{}
			translate.Language = e.ChildText("div.langsmall-song span.mobile-only-inline")
			translate.Author.Name = e.ChildText(".authorsubmitted a")
			translate.Author.Href = e.ChildAttr(".authorsubmitted a[href]", "href")
			if translate.Language != "" {
				translate.Title = e.ChildText("h2.title-h2")
				e.ForEach(".ltf .par div, .emptyline", func(_ int, e *colly.HTMLElement) {
					translate.Lines = append(translate.Lines, e.Text)
				})
				l.Translates = append(l.Translates, translate)
			}
		})

		// TODO
		// Fix more-then-one translate issue.
		if strings.Contains(allowedTranslationLanguages, e.Text) {
			cc.Visit("https://lyricstranslate.com/" + e.Attr("href"))
		}
	})

	c.Visit(url)
}
