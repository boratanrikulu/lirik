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
		} `json:"items"`
	} `json:"tracks"`
}

func (s *SpotifyAuthorizatied) Search(search string) (Song, error) {
	url, err := url.Parse("https://api.spotify.com/v1/search")
	if err != nil {
		log.Print(err)
		return Song{}, errors.New("Error.")
	}

	quaries := url.Query()
	quaries.Add("q", search)
	quaries.Add("type", "track")
	quaries.Add("limit", "1")
	url.RawQuery = quaries.Encode()
	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		log.Print(err)
		return Song{}, errors.New("Error.")
	}
	req.Header.Set("Authorization", "Bearer "+s.Authorization.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Song{}, errors.New("Error.")
	}
	defer resp.Body.Close()
	if resp.StatusCode > 399 {
		log.Print(resp.StatusCode)
		return Song{}, errors.New("Error.")
	}

	var response searchResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if len(response.Tracks.Items) == 0 {
		return Song{}, errors.New("We can not find any related song.")
	}

	var searchSong Song
	searchSong.ArtistName = response.Tracks.Items[0].Artists[0].Name
	searchSong.Name = response.Tracks.Items[0].Name
	searchSong.AlbumName = response.Tracks.Items[0].Album.Name
	if response.Tracks.Items[0].Album.AlbumType == "single" {
		searchSong.AlbumName += " (Single)"
	}
	searchSong.ReleaseDate = response.Tracks.Items[0].Album.ReleaseDate
	searchSong.TotalTracks = response.Tracks.Items[0].Album.TotalTracks
	albumImages := response.Tracks.Items[0].Album.Images
	if len(albumImages) >= 2 {
		searchSong.AlbumImage = albumImages[len(albumImages)-2].URL
	}

	return searchSong, nil
}
