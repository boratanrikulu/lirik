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
	RequestAuthorization RequestAuthorization
}

// Public Methods

func (s Spotify) GetSpotifyAuthLink() string {
	SetRequestAuthorization(&s.RequestAuthorization)

	baseUrl, _ := url.Parse(s.RequestAuthorization.BaseURL)
	baseUrl.Path += s.RequestAuthorization.Path

	params := url.Values{}
	params.Add("client_id", s.RequestAuthorization.ClientID)
	params.Add("response_type", s.RequestAuthorization.ResponseType)
	params.Add("redirect_uri", s.RequestAuthorization.RedirectURI)
	params.Add("scope", s.RequestAuthorization.Scope)
	params.Add("state", s.RequestAuthorization.State)
	baseUrl.RawQuery = params.Encode()
	s.RequestAuthorization.Link = baseUrl.String()

	return s.RequestAuthorization.Link
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
