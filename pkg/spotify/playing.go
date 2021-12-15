package spotify

import (
	"encoding/json"
	"errors"
	"net/http"
)

type currentlyPlaying struct {
	Item struct {
		Album struct {
			Name        string `json:"name"`
			AlbumType   string `json:"album_type"`
			ReleaseDate string `json:"release_date"`
			TotalTracks int    `json:"total_tracks"`
			Images      []struct {
				URL string `json:"url"`
			} `json:"images"`
		} `json:"album"`
		Artists []struct {
			Name string `json:"name"`
		} `json:"artists"`
		Name string `json:"name"`
	} `json:"item"`
}

func (s *SpotifyAuthorizatied) GetCurrentlyPlaying() (Song, error) {
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me/player/currently-playing", nil)
	req.Header.Set("Authorization", "Bearer "+s.Authorization.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Song{}, errors.New("Error.")
	}
	defer resp.Body.Close()

	var response currentlyPlaying
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil || response.Item.Name == "" {
		return Song{}, errors.New("Error.")
	}

	var currentSong Song
	currentSong.Name = response.Item.Name
	currentSong.ArtistName = response.Item.Artists[0].Name
	currentSong.AlbumName = response.Item.Album.Name
	if response.Item.Album.AlbumType == "single" {
		currentSong.AlbumName += " (Single)"
	}
	currentSong.TotalTracks = response.Item.Album.TotalTracks
	currentSong.ReleaseDate = response.Item.Album.ReleaseDate
	albumImages := response.Item.Album.Images
	if len(albumImages) >= 2 {
		currentSong.AlbumImage = albumImages[len(albumImages)-2].URL
	}

	return currentSong, nil
}
