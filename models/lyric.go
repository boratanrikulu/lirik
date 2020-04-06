package models

import (
	"encoding/json"
	"errors"

	//"bytes"
	//"bytes"
	//"encoding/json"
	//"errors"
	"regexp"
	"fmt"
	"io/ioutil"
	"net/url"
	//"strings"

	//"net/url"

	//"io/ioutil"
	"net/http"
	//"net/url"
)

type Lyric struct {
	SongName   string
	ArtistName string
	GeniusID   int
}

type Genius struct {
	Response struct {
		Hits []struct {
			Type string `json:type`
			Result struct {
				ID int `json:"id"`
			} `json:"result"`
		} `json:"hits"`
	} `json:"response"`
}

// Public Methods

func (l *Lyric) GetLyric() error {
	// Removes "(....)" values from song name.
	re := regexp.MustCompile(`\(.*?\)`)
	songName := re.ReplaceAllString(l.SongName, "")

	u, _ := url.Parse("https://api.genius.com/search")
	q, _ := url.ParseQuery(u.RawQuery)

	q.Add("q", l.ArtistName + " " + songName)
	u.RawQuery = q.Encode()
	req, _ := http.NewRequest("GET", fmt.Sprint(u), nil)

	req.Header.Set("Authorization",
		"Bearer IwH8C9cJrsFo3rMg86h9kQ0DWP04ytVSwIh8906uSOWKQS7aDPCLdNlZEB7xDTwx")

	// Sends the request.
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	//fmt.Println(params.Encode())
	//fmt.Println(string(body))
	// Reads response and unmarshal it to spotify model.
	//fmt.Println(resp)
	body, _ := ioutil.ReadAll(resp.Body)
	genius := new(Genius)
	json.Unmarshal(body, &genius)
	//fmt.Println(string(body))
	//fmt.Println(genius.Response)
	fmt.Println(l.SongName, l.ArtistName)
	if len(genius.Response.Hits) != 0 {
		for _, value := range genius.Response.Hits {
			if value.Type == "song" {
				fmt.Println(value.Type, value.Result.ID)
				l.GeniusID = value.Result.ID
				return nil
			}
		}
	}
	return errors.New("Error.")
}
