package models

import (
	"crypto/rand"
	"fmt"
	"net/url"
)

type RequestAuthorization struct {
	Link string
	BaseURL string
	Path string
	ClientID string
	ResponseType string
	RedirectURI string
	Scope string
	State string
}

type Spotify struct {
}

// Public Methods

func (s Spotify) GetSpotifyAuthLink() string {
	var requestAuthorization RequestAuthorization
	SetRequestAuthorization(&requestAuthorization)

	baseUrl, _ := url.Parse(requestAuthorization.BaseURL)
	baseUrl.Path += requestAuthorization.Path

	params := url.Values{}
	params.Add("client_id", requestAuthorization.ClientID)
	params.Add("response_type", requestAuthorization.ResponseType)
	params.Add("redirect_uri", requestAuthorization.RedirectURI)
	params.Add("scope", requestAuthorization.Scope)
	params.Add("state", requestAuthorization.State)
	baseUrl.RawQuery = params.Encode()
	requestAuthorization.Link = baseUrl.String()

	return requestAuthorization.Link
}

// Private Methods

func SetRequestAuthorization(requestAuthorization *RequestAuthorization) {
	requestAuthorization.BaseURL      = "https://accounts.spotify.com/"
	requestAuthorization.Path         = "authorize"
	requestAuthorization.ClientID     = "6f524a004e874120b42251c6c6d0e699"
	requestAuthorization.ResponseType = "code"
	requestAuthorization.RedirectURI  = "http://localhost:3000/spotify"
	requestAuthorization.Scope        = "user-read-currently-playing streaming user-read-playback-state"
	requestAuthorization.State        = randomState()
}

func randomState() string {
	key := make([]byte, 32)
	rand.Read(key)
	return fmt.Sprintf("%x", key)
}
