package controllers

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/boratanrikulu/lirik.app/controllers/helpers"
	"github.com/boratanrikulu/lirik.app/models"
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
	if accessTokenCookie == nil && refreshTokenCookie == nil {
		// Cookie is not exist.
		// That means user does not have tokens.
		// We will send a request with code value to take tokens.
		err := takeTokens(spotify, w, r)
		if err != nil {
			return
		}
		// Else,
		// Redirect to welcome page.
		// To do not show spotify's callback query.
		// For just cosmetic.
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else if accessTokenCookie == nil && refreshTokenCookie != nil {
		// That means we have refresh token,
		// but our access token has expired.
		// We need to send a request to spotify
		// to take new access token.
		// We will use refresh token to take access token.
		// Response will be include access token.
		// Also, response might be include refresh token.
		// If response has refresh token, we will update our refresh token.

		// Set refresh token to spotify object.
		tokenResponse := &spotify.RefreshAndAccessTokens.Response
		tokenResponse.RefreshToken = refreshTokenCookie.Value

		err := updateTokens(spotify, w, r)
		if err != nil {
			return
		}
		// Else,
		// Redirect to welcome page.
		// We couldn't get access token.
		// TODO show an error message on welcome page.
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		// Cookie is exist.
		// Set it to the spotify object.
		tokenResponse := &spotify.RefreshAndAccessTokens.Response
		tokenResponse.AccessToken = accessTokenCookie.Value
		// Set refresh token if it is not empty.
		// We do not need to refresh token to showing song.
		// So,
		// If there is only access token and not have refresh token,
		// it is not a problem.
		// That means user will be login for access token's expire time.
		if refreshTokenCookie != nil {
			tokenResponse.RefreshToken = refreshTokenCookie.Value
		}
	}

	// Gets user's current song.
	artistName, songName, albumImage, err := spotify.GetCurrentlyPlaying()
	if err != nil {
		errorMessages := []string{
			"There is no song playing.",
			"You need to play a song! 😅",
			"",
			"Open your spotify account and play a song. 🎶 🎉",
		}
		helpers.ErrorPage(errorMessages, w)
		return
	}

	// Show lyrics result.
	showLyric(artistName, songName, albumImage, w, r)
}

// Private Methods

func showLyric(artistName string, songName string, albumImage string, w http.ResponseWriter, r *http.Request) {
	// Set params.
	q, _ := url.ParseQuery("")
	q.Add("artistName", artistName)
	q.Add("songName", songName)
	q.Add("albumImage", albumImage)

	// Update request with created url.
	r.URL.RawQuery = q.Encode()

	// Handle request with lyric controller.
	LyricGet(w, r)
}

func updateTokens(spotify *models.Spotify, w http.ResponseWriter, r *http.Request) error {
	err := spotify.GetUpdateTokens()
	if err != nil {
		return err
	}
	// Everthinng is okay,
	// Set access token to cookies.
	// Also, it will update access token if response has a new access token.
	helpers.UpdateTokenCookies(spotify.UpdateAccessToken, w)
	return nil
}

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
		// Redirect back to "/".
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return fmt.Errorf("Error.", nil)
	}

	// Take state cookie to check client and the response are the same.
	state, err := r.Cookie("State")
	if err != nil {
		// If there is no state cookie that means user has deleted it.
		// Show error and do not anything.
		log.Printf("[ERROR] Error occur. User deleted it's state cookie: %v", err)
		errorMessages := []string{
			"There is no state cookie.",
			"Please do not remove your state cookie.",
		}
		helpers.ErrorPage(errorMessages, w)
		return err
	}

	if state.Value != authResponse.State {
		// If it not same that might be an attack.
		// Show the error message and do not anything.
		log.Print("[ERROR] User's state and request state are not same.")
		errorMessages := []string{
			"Your state cookie and the response are not same.",
			"You might be under attack.",
		}
		helpers.ErrorPage(errorMessages, w)
		return err
	}

	// State valeus are equal.
	// It is safe to make request to take user's tokens.
	// Send the request. Get response. Set values to spotify object.
	err = spotify.GetRefreshAndAccessTokensResponse()
	if err != nil {
		// If there is a error, log it.
		// Say to user, we have some issues.
		log.Printf("[ERROR] Error occur while taking response from spotify: %v", err)
		errorMessages := []string{
			"There is some issues while taking response from Spotify.",
			"Please try later.",
		}
		helpers.ErrorPage(errorMessages, w)
	}
	// If everything is okay,
	// Then set the tokens to cookie for later usage..
	helpers.SetTokenCookies(spotify.RefreshAndAccessTokens, w)
	return nil

}
