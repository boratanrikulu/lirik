package spotify

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

func (s *Spotify) Login(state, code string) (sAuthed SpotifyAuthorizatied, err error) {
	if state == "" || code == "" {
		return sAuthed, errors.New("state or code can not be empty")
	}

	response, err := s.getTokens(state, code)
	if err != nil {
		return sAuthed, err
	}

	sAuthed.Authorization = response
	sAuthed.UserMe = sAuthed.GetUserMe()

	return sAuthed, nil
}

func (s *Spotify) getTokens(state, code string) (response AuthorizationResponse, err error) {
	params := url.Values{}
	params.Add("code", code)
	params.Add("grant_type", "authorization_code")
	params.Add("redirect_uri", s.RedirectURI)

	req, _ := http.NewRequest("POST", "https://accounts.spotify.com/api/token",
		bytes.NewBuffer([]byte(params.Encode())))
	req.Header.Set("Authorization", S.Authorization)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return response, err
	}

	return response, nil
}

type SpotifyAuthorizatied struct {
	UserMe        string
	Authorization AuthorizationResponse
}

type AuthorizationResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

func (s SpotifyAuthorizatied) GetUpdatedTokens() (response AuthorizationResponse, err error) {
	params := url.Values{}
	params.Add("grant_type", "refresh_token")
	params.Add("refresh_token", s.Authorization.RefreshToken)

	req, _ := http.NewRequest("POST", "https://accounts.spotify.com/api/token",
		bytes.NewBuffer([]byte(params.Encode())))
	req.Header.Set("Authorization", S.Authorization)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return response, err
	}

	return response, nil
}

func (s SpotifyAuthorizatied) GetUserMe() string {
	req, _ := http.NewRequest("GET", "https://api.spotify.com/v1/me", nil)
	req.Header.Set("Authorization", "Bearer "+s.Authorization.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		log.Println(resp.StatusCode)
		return ""
	}

	spotifyUserMe := make(map[string]interface{})
	err = json.NewDecoder(resp.Body).Decode(&spotifyUserMe)
	if err != nil {
		log.Println(err)
		return ""
	}

	userMe, ok := spotifyUserMe["uri"].(string)
	if !ok {
		return ""
	}

	return userMe
}

func randomState() string {
	key := make([]byte, 32)
	rand.Read(key)
	return fmt.Sprintf("%x", key)
}
