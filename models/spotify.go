package models

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
)

type RequestAuthorization struct {
	Link string
	BaseURL string
	Path string
	ResponseType string
	Scope string
	State string
}

type RefreshAndAccessTokens struct {
	URL string
	Code string
	GrantType string
	RedirectURI string
	Authorization string
	ContentType string
}

type Spotify struct {
	ClientID string
	ClientSecret string
	RedirectURI string
	RequestAuthorization RequestAuthorization
	RefreshAndAccessTokens RefreshAndAccessTokens
}

// Public Methods

func (s *Spotify) InitSecrets() {
	s.ClientID     = "6f524a004e874120b42251c6c6d0e699"
	s.ClientSecret = "2ed3ffdd211a4f2ab38d6da112316fee"
	s.RedirectURI  = "http://localhost:3000/spotify"
}

func (s *Spotify) GetRequestAuthorizationLink() string {
	setRequestAuthorization(&s.RequestAuthorization)

	baseUrl, _ := url.Parse(s.RequestAuthorization.BaseURL)
	baseUrl.Path += s.RequestAuthorization.Path

	params := url.Values{}
	params.Add("client_id", s.ClientID)
	params.Add("response_type", s.RequestAuthorization.ResponseType)
	params.Add("redirect_uri", s.RedirectURI)
	params.Add("scope", s.RequestAuthorization.Scope)
	params.Add("state", s.RequestAuthorization.State)
	baseUrl.RawQuery = params.Encode()
	s.RequestAuthorization.Link = baseUrl.String()

	return s.RequestAuthorization.Link
}

func (s *Spotify) GetRefreshAndAccessTokensReq(code string) *http.Request {
	setRefreshAndAccessTokens(s, code)
	r := s.RefreshAndAccessTokens

	params := url.Values{}
	params.Add("code", r.Code)
	params.Add("grant_type", r.GrantType)
	params.Add("redirect_uri", s.RedirectURI)

	req, _ := http.NewRequest("POST", r.URL,
		bytes.NewBuffer([]byte(params.Encode())))
	req.Header.Set("Authorization", r.Authorization)
	req.Header.Set("Content-Type", r.ContentType)

	return req
}

// Private Methods

func setRequestAuthorization(r *RequestAuthorization) {
	r.BaseURL      = "https://accounts.spotify.com/"
	r.Path         = "authorize"
	r.ResponseType = "code"
	r.Scope        = "user-read-currently-playing streaming user-read-playback-state"
	r.State        = randomState()
}

func setRefreshAndAccessTokens(s *Spotify, code string) {
	secrets := s.ClientID + ":" + s.ClientSecret
	encoded := "Basic " + base64.StdEncoding.EncodeToString([]byte(secrets))

	s.RefreshAndAccessTokens.URL           = "https://accounts.spotify.com/api/token"
	s.RefreshAndAccessTokens.Code          = code
	s.RefreshAndAccessTokens.GrantType     = "authorization_code"
	s.RefreshAndAccessTokens.RedirectURI   = "http://localhost:3000/spotify"
	s.RefreshAndAccessTokens.Authorization = encoded
	s.RefreshAndAccessTokens.ContentType   = "application/x-www-form-urlencoded"
}

func randomState() string {
	key := make([]byte, 32)
	rand.Read(key)
	return fmt.Sprintf("%x", key)
}
