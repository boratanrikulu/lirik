package spotify

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
)

type searchResponse struct {
	Tracks struct {
		Items []struct {
			Album struct {
				Images []struct {
					URL string `json:"url"`
				} `json:"images"`
			} `json:"album"`
			Artists []struct {
				Name string `json:"name"`
			} `json:"artists"`
			Name string `json:"name"`
		} `json:"items"`
	} `json:"tracks"`
}

func (s *SpotifyAuthorizatied) Search(search string) (artistName string, songName string, albumImage string, err error) {
	url, err := url.Parse("https://api.spotify.com/v1/search")
	if err != nil {
		log.Print(err)
		return "", "", "", errors.New("Error.")
	}

	quaries := url.Query()
	quaries.Add("q", search)
	quaries.Add("type", "track")
	quaries.Add("limit", "1")
	url.RawQuery = quaries.Encode()
	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		log.Print(err)
		return "", "", "", errors.New("Error.")
	}
	req.Header.Set("Authorization", "Bearer "+s.Authorization.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", "", errors.New("Error.")
	}
	defer resp.Body.Close()
	if resp.StatusCode > 399 {
		log.Print(resp.StatusCode)
		return "", "", "", errors.New("Error.")
	}

	var response searchResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if len(response.Tracks.Items) == 0 {
		return "", "", "", errors.New("We can not find any related song.")
	}

	artistName = response.Tracks.Items[0].Artists[0].Name
	songName = response.Tracks.Items[0].Name
	albumImages := response.Tracks.Items[0].Album.Images
	if len(albumImages) >= 2 {
		albumImage = albumImages[len(albumImages)-2].URL
	}

	return artistName, songName, albumImage, nil
}
