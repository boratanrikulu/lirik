package models

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

type Spotify struct {
	ClientID               string
	ClientSecret           string
	RedirectURI            string
	Authorization          Authorization
	RefreshAndAccessTokens RefreshAndAccessTokens
	CurrentlyPlaying       CurrentlyPlaying
}

type Authorization struct {
	Request struct {
		Link         string
		BaseURL      string
		Path         string
		ResponseType string
		Scope        string
		State        string
	}
	Response struct {
		Code		string
		State		string
	}
}

type RefreshAndAccessTokens struct {
	Request struct {
		URL           string
		Code          string
		GrantType     string
		RedirectURI   string
		Authorization string
		ContentType   string
	}
	Response struct {
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
		Scope        string `json:"scope"`
	}
}

type CurrentlyPlaying struct {
	Item struct {
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
}

// Public Methods

func (s *Spotify) InitSecrets() {
	s.ClientID = os.Getenv("CLIENT_ID")
	s.ClientSecret = os.Getenv("CLIENT_SECRET")
	s.RedirectURI = os.Getenv("REDIRECT_URI")

	if s.ClientID == "" || s.ClientSecret == "" || s.RedirectURI == "" {
		panic("Secrets are not set in env file.")
	}
}

func (s *Spotify) GetRequestAuthorizationLink() (authLink string, err error) {
	setRequestAuthorization(s)
	r := &s.Authorization.Request

	baseUrl, _ := url.Parse(r.BaseURL)
	baseUrl.Path += r.Path

	params := url.Values{}
	params.Add("client_id", s.ClientID)
	params.Add("response_type", r.ResponseType)
	params.Add("redirect_uri", s.RedirectURI)
	params.Add("scope", r.Scope)
	params.Add("state", r.State)
	baseUrl.RawQuery = params.Encode()
	r.Link = baseUrl.String()

	return r.Link, nil
}

func (s *Spotify) GetRefreshAndAccessTokensResponse() error {
	setRefreshAndAccessTokens(s)
	r := s.RefreshAndAccessTokens.Request

	params := url.Values{}
	params.Add("code", r.Code)
	params.Add("grant_type", r.GrantType)
	params.Add("redirect_uri", s.RedirectURI)

	req, _ := http.NewRequest("POST", r.URL,
		bytes.NewBuffer([]byte(params.Encode())))
	req.Header.Set("Authorization", r.Authorization)
	req.Header.Set("Content-Type", r.ContentType)

	// Sends the request.
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Reads response and unmarshal it to spotify model.
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &s.RefreshAndAccessTokens.Response)
	if err != nil {
		return err
	}
	return nil
}

func (s *Spotify) GetCurrentlyPlaying() (artistName string, songName string, err error) {
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me/player/currently-playing", nil)
	req.Header.Set("Authorization", "Bearer "+s.RefreshAndAccessTokens.Response.AccessToken)
	fmt.Println(s.RefreshAndAccessTokens.Response.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", errors.New("Error.")
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &s.CurrentlyPlaying)
	if s.CurrentlyPlaying.Item.Name == "" {
		return "", "", errors.New("Error.")
	}

	return s.CurrentlyPlaying.Item.Artists[0].Name, s.CurrentlyPlaying.Item.Name, nil
}

// Private Methods

func setRequestAuthorization(s *Spotify) {
	r := &s.Authorization.Request
	r.BaseURL = "https://accounts.spotify.com/"
	r.Path = "authorize"
	r.ResponseType = "code"
	r.Scope = "user-read-currently-playing streaming user-read-playback-state"
	r.State = randomState()
}

func setRefreshAndAccessTokens(s *Spotify) {
	secrets := s.ClientID + ":" + s.ClientSecret
	encoded := "Basic " + base64.StdEncoding.EncodeToString([]byte(secrets))

	r := &s.RefreshAndAccessTokens.Request
	r.URL = "https://accounts.spotify.com/api/token"
	r.Code = s.Authorization.Response.Code
	r.GrantType = "authorization_code"
	r.RedirectURI = "http://localhost:3000/spotify"
	r.Authorization = encoded
	r.ContentType = "application/x-www-form-urlencoded"
}

func randomState() string {
	key := make([]byte, 32)
	rand.Read(key)
	return fmt.Sprintf("%x", key)
}
