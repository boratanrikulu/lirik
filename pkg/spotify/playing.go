package spotify

import (
	"encoding/json"
	"errors"
	"net/http"
)

type currentlyPlaying struct {
	Item struct {
		Album struct {
			Images []struct {
				URL string `json:"url"`
			} `json:"images"`
		} `json:"album"`
		Artists []struct {
			Name string `json:"name"`
		} `json:"artists"`
		Name string `json:"name"`
	} `json:"item"`
}

func (s *SpotifyAuthorizatied) GetCurrentlyPlaying() (artistName string, songName string, albumImage string, err error) {
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me/player/currently-playing", nil)
	req.Header.Set("Authorization", "Bearer "+s.Authorization.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", "", errors.New("Error.")
	}
	defer resp.Body.Close()

	var response currentlyPlaying
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil || response.Item.Name == "" {
		return "", "", "", errors.New("Error.")
	}

	artistName = response.Item.Artists[0].Name
	songName = response.Item.Name
	albumImages := response.Item.Album.Images
	if len(albumImages) >= 2 {
		albumImage = albumImages[len(albumImages)-2].URL
	}

	return artistName, songName, albumImage, nil
}
