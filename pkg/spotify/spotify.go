package spotify

import (
	"encoding/base64"
	"net/url"
)

type Spotify struct {
	ClientID      string
	Authorization string
	RedirectURI   string
}

var S *Spotify

func NewSpotify(clientID, clientSecret, redirectURI string) *Spotify {
	secrets := clientID + ":" + clientSecret
	encoded := "Basic " + base64.StdEncoding.EncodeToString([]byte(secrets))

	return &Spotify{
		ClientID:      clientID,
		Authorization: encoded,
		RedirectURI:   redirectURI,
	}
}

func (s *Spotify) GetRequestAuthorizationLink() (authLink, state string, err error) {
	baseUrl, err := url.Parse("https://accounts.spotify.com/authorize")
	if err != nil {
		return authLink, state, err
	}

	params := url.Values{}
	params.Add("client_id", s.ClientID)
	params.Add("response_type", "code")
	params.Add("redirect_uri", s.RedirectURI)
	params.Add("scope", "user-read-currently-playing")
	state = randomState()
	params.Add("state", state)
	baseUrl.RawQuery = params.Encode()

	return baseUrl.String(), state, nil
}
