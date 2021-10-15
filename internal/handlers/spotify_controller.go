package handlers

import (
	"log"
	"net/http"
	"net/url"

	"github.com/boratanrikulu/lirik.app/internal/handlers/helpers"
	"github.com/boratanrikulu/lirik.app/pkg/spotify"
)

func SpotifyGet(w http.ResponseWriter, r *http.Request) {
	accessTokenCookie, _ := r.Cookie("AccessToken")
	refreshTokenCookie, _ := r.Cookie("RefreshToken")
	stateCookie, _ := r.Cookie("State")
	meCookie, _ := r.Cookie("Me")

	var sAuthed spotify.SpotifyAuthorizatied
	var err error
	if accessTokenCookie == nil && refreshTokenCookie == nil {
		state := r.URL.Query().Get("state")
		if stateCookie == nil || stateCookie.Value != state {
			log.Print("[ERROR] User's state and request state are not same.")
			errorMessages := []string{
				"Your state cookie and the response are not same.",
				"You might be under attack.",
			}

			helpers.ErrorPage(errorMessages, w)
			return
		}
		code := r.URL.Query().Get("code")
		sAuthed, err = spotify.S.Login(state, code)
		if err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}

		helpers.UpdateTokenCookies(sAuthed, w)
		http.Redirect(w, r, "/spotify", http.StatusSeeOther)
		return
	}

	if accessTokenCookie == nil && refreshTokenCookie != nil {
		sAuthed.Authorization.RefreshToken = refreshTokenCookie.Value
		sAuthed.Authorization, err = sAuthed.GetUpdatedTokens()
		if err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		sAuthed.UserMe = sAuthed.GetUserMe()

		helpers.UpdateTokenCookies(sAuthed, w)
	}

	if accessTokenCookie != nil {
		sAuthed.Authorization.AccessToken = accessTokenCookie.Value
	}
	if refreshTokenCookie != nil {
		sAuthed.Authorization.RefreshToken = refreshTokenCookie.Value
	}
	if meCookie != nil {
		sAuthed.UserMe = meCookie.Value
	}

	log.Printf("[USER] %s\n", sAuthed.UserMe)
	artistName, songName, albumImage, err := sAuthed.GetCurrentlyPlaying()
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

	showLyric(artistName, songName, albumImage, w, r)
}

func showLyric(artistName string, songName string, albumImage string, w http.ResponseWriter, r *http.Request) {
	q, _ := url.ParseQuery("")
	q.Add("artistName", artistName)
	q.Add("songName", songName)
	q.Add("albumImage", albumImage)
	r.URL.RawQuery = q.Encode()
	LyricsGet(w, r)
}
