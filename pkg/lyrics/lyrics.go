package lyrics

import (
	"log"
	"regexp"
	"strings"
)

// Finder finds lyrics by supporting many sources usage.
type Finder interface {
	GetLyrics(string, string) (bool, Lyrics)
}

type finder struct {
	sources []source
}

type source struct {
	name   string
	maxTry int
	f      Finder
}

func NewFinder() Finder {
	return &finder{
		sources: []source{
			{
				name:   "LOCAL",
				maxTry: 1,
				f:      newLocalSource(),
			},
			{
				name:   "FIRST",
				maxTry: 1,
				f:      newFirstSource(),
			},
			{
				name:   "THIRD",
				maxTry: 4,
				f:      newSecondSource(),
			},
		},
	}
}

// GetLyrics returns related lyrics for the song.
func (f *finder) GetLyrics(artistName string, songName string) (found bool, lyrics Lyrics) {
	songName = songRegex(songName)

	for _, source := range f.sources {
		for i := 0; i < source.maxTry; i++ {
			if f, lyrics := source.f.GetLyrics(artistName, songName); f {
				if source.name != "LOCAL" {
					go saveToFile(artistName, songName, lyrics)
				}
				log.Printf("[%s] [FOUND] [x%d] \"%s by %s\"", source.name, i+1, songName, artistName)
				return true, lyrics
			}
		}
	}

	log.Printf("[ALL] [NOT FOUND] \"%s by %s\"", songName, artistName)
	return false, lyrics
}

type Lyrics struct {
	Lines      []string
	Language   string
	Translates []Translate `json:"-"`
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

func songRegex(song string) string {
	regexList := []string{
		` - .+`,                     // Removes values after " - ...". from song name.
		`(?i)\(.*?feat.*?\)`,        // Removes all (...feat...)s from song name.
		`(?i)\[.*?feat.*?\]`,        // Removes all [...feat..]s from song name.
		`(?i)\(.*?remastered.*?\)`,  // Removes all (...remastered...)s from song name.
		`(?i)\[.*?remastered.*?\)]`, // Removes all [...remastered...]s from song name.
		`(?i)\(.*?cover.*?\)`,       // Removes all (...cover...)s from song name.
		`(?i)\[.*?cover.*?\]`,       // Removes all [...cover...]s from song name.
		`(?i)\(.*?with.*?\)`,        // Removes all (...with...)s from song name.
		`(?i)\[.*?with.*?\]`,        // Removes all [...with...]s from song name.
	}

	for _, value := range regexList {
		re := regexp.MustCompile(value)
		song = re.ReplaceAllString(song, "")
	}

	song = strings.TrimSpace(song)
	return song
}
