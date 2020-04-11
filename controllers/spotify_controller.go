package controllers

import (
	"io/ioutil"
	"net/url"
	"log"
	"fmt"
	"github.com/boratanrikulu/s-lyrics/models"
	"github.com/boratanrikulu/s-lyrics/controllers/helpers"
	"net/http"
)

// Public Methods

func SpotifyGet(w http.ResponseWriter, r *http.Request) {	
	// Creates a spotify model with it's secrets.
	spotify := new(models.Spotify)
	spotify.InitSecrets()

	// Get tokens from cookie.
	accessTokenCookie, _ := r.Cookie("AccessToken")
	refreshTokenCookie, _ := r.Cookie("RefreshToken")

	// Check tokens situation.
	if accessTokenCookie == nil || refreshTokenCookie == nil  {
		// Cookie is not exist.
		// That means user does not have tokens.
		// We will send a request with code value to take tokens.
		err := takeTokens(spotify, w, r)
		if err != nil {
			return
		}
	} else {
		// Cookie is exist.
		// Set it to the spotify object.
		tokenResponse := &spotify.RefreshAndAccessTokens.Response
		tokenResponse.AccessToken = accessTokenCookie.Value
		tokenResponse.RefreshToken = refreshTokenCookie.Value
	}

	// Gets user's current song.
	artistName, songName, err := spotify.GetCurrentlyPlaying()
	if err != nil {
		helpers.ErrorPage("Çalınan bir şarkı yok.", w)
		return
	}

	// Show lyrics result.
	showLyric(artistName, songName, w)
}

// Private Methods

func takeTokens(spotify *models.Spotify, w http.ResponseWriter, r *http.Request) error {
	// That means user does not have tokens.
	// We will send a request with code value to take tokens.
	// But there is some important things to check first.

	// If user does not have tokens,
	// Then there must be a token value on the request.
	// Get the code and state from request query.
	authResponse := &spotify.Authorization.Response
	authResponse.State = r.URL.Query().Get("state")
	authResponse.Code = r.URL.Query().Get("code")
	if authResponse.Code == "" {
		// That means user did not auth him/her spotify account,
		// and trying to open /spotify path.
		// Show an error and do not anything.
		errorMessage := "Seems like you do not auth your spotify account."
		errorMessage += "\nPlease go to welcome page and click to \"Login with Spotify Account\" button."
		helpers.ErrorPage(errorMessage, w)
		return fmt.Errorf("Error.", nil)
	}

	// Take state cookie to check client and the response are the same.
	state, err := r.Cookie("State")
	if err != nil {
		// If there is no state cookie that means user has deleted it.
		// Show error and do not anything.
		log.Print("Error occur. User deleted it's state cookie: %v", err)
		errorMessage := "There is no state cookie."
		errorMessage += "\nPlease do not remove your state cookie."
		helpers.ErrorPage(errorMessage, w)
		return err
	}

	if state.Value != authResponse.State {
		// If it not same that might be an attack.
		// Show the error message and do not anything.
		log.Print("User's state and request state are not same.")
		errorMessage := "Your state cookie and the response are not same."
		errorMessage += "\nYou might be under attack."
		helpers.ErrorPage(errorMessage, w)
		return err
	}

	// State valeus are equal.
	// It is safe to make request to take user's tokens.
	// Send the request. Get response. Set values to spotify object.
	err = spotify.GetRefreshAndAccessTokensResponse()
	if err != nil {
		// If there is a error, log it.
		// Say to user, we have some issues.
		log.Print("Error occur while taking response from spotify: %v", err)
		errorMessage := "There is some issues while taking response from Spotify."
		errorMessage += "\nPlease try later."
		helpers.ErrorPage(errorMessage, w)
	}
	// If everything is okay,
	// Then set the tokens to cookie for later usage..
	helpers.SetTokenCookies(spotify.RefreshAndAccessTokens, w)
	return nil
}

func showLyric(artistName string, songName string, w http.ResponseWriter) {
	// Set params.
	u, _ := url.Parse("/lyric")
	q, _ := url.ParseQuery(u.RawQuery)
	q.Add("artistName", artistName)
	q.Add("songName", songName)
	u.RawQuery = q.Encode()

	// Send request to lyric_controller.
	// TODO: fix this url.
	urlQuery := "http://localhost:3000" + fmt.Sprint(u)
	resp, _ := http.Get(urlQuery)
	defer resp.Body.Close()

	// Show the result from lyric_controller.
	body, _ := ioutil.ReadAll(resp.Body)
	w.Write(body)
}
