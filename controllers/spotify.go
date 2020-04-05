package controllers

import (
	// "strings"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

type CurrentlyPlaying struct {
	Timestamp int64 `json:"timestamp"`
	Context   struct {
		ExternalUrls struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`
		Href string `json:"href"`
		Type string `json:"type"`
		URI  string `json:"uri"`
	} `json:"context"`
	ProgressMs int `json:"progress_ms"`
	Item       struct {
		Album struct {
			AlbumType string `json:"album_type"`
			Artists   []struct {
				ExternalUrls struct {
					Spotify string `json:"spotify"`
				} `json:"external_urls"`
				Href string `json:"href"`
				ID   string `json:"id"`
				Name string `json:"name"`
				Type string `json:"type"`
				URI  string `json:"uri"`
			} `json:"artists"`
			AvailableMarkets []string `json:"available_markets"`
			ExternalUrls     struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Href   string `json:"href"`
			ID     string `json:"id"`
			Images []struct {
				Height int    `json:"height"`
				URL    string `json:"url"`
				Width  int    `json:"width"`
			} `json:"images"`
			Name                 string `json:"name"`
			ReleaseDate          string `json:"release_date"`
			ReleaseDatePrecision string `json:"release_date_precision"`
			TotalTracks          int    `json:"total_tracks"`
			Type                 string `json:"type"`
			URI                  string `json:"uri"`
		} `json:"album"`
		Artists []struct {
			ExternalUrls struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Href string `json:"href"`
			ID   string `json:"id"`
			Name string `json:"name"`
			Type string `json:"type"`
			URI  string `json:"uri"`
		} `json:"artists"`
		AvailableMarkets []string `json:"available_markets"`
		DiscNumber       int      `json:"disc_number"`
		DurationMs       int      `json:"duration_ms"`
		Explicit         bool     `json:"explicit"`
		ExternalIds      struct {
			Isrc string `json:"isrc"`
		} `json:"external_ids"`
		ExternalUrls struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`
		Href        string `json:"href"`
		ID          string `json:"id"`
		IsLocal     bool   `json:"is_local"`
		Name        string `json:"name"`
		Popularity  int    `json:"popularity"`
		PreviewURL  string `json:"preview_url"`
		TrackNumber int    `json:"track_number"`
		Type        string `json:"type"`
		URI         string `json:"uri"`
	} `json:"item"`
	CurrentlyPlayingType string `json:"currently_playing_type"`
	Actions              struct {
		Disallows struct {
			Resuming bool `json:"resuming"`
		} `json:"disallows"`
	} `json:"actions"`
	IsPlaying bool `json:"is_playing"`
}

type JSONData struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

func SpotifyGet(w http.ResponseWriter, r *http.Request) {
	params := url.Values{}
	params.Add("code", r.URL.Query().Get("code"))
	params.Add("grant_type", "authorization_code")
	params.Add("redirect_uri", "http://localhost:3000/spotify")

	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token",
		bytes.NewBuffer([]byte(params.Encode())))

	msg := "6f524a004e874120b42251c6c6d0e699:2ed3ffdd211a4f2ab38d6da112316fee"
	encoded := "Basic " + base64.StdEncoding.EncodeToString([]byte(msg))
	req.Header.Set("Authorization", encoded)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	jsonData := JSONData{}
	json.Unmarshal(body, &jsonData)

	currentSong, currentArtist := currentlyPlaying(jsonData.AccessToken)

	a := "?songName=" + currentSong + "&artistName=" + currentArtist
	http.Redirect(w, r, "/a"+a, 300)
}

func currentlyPlaying(access_token string) (string, string) {
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me/player/currently-playing", nil)
	req.Header.Set("Authorization", "Bearer "+access_token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	currentlyPlaying := CurrentlyPlaying{}
	json.Unmarshal(body, &currentlyPlaying)

	return currentlyPlaying.Item.Name, currentlyPlaying.Item.Artists[0].Name
}
