package controllers

import (
	"github.com/boratanrikulu/s-lyrics/models"
	"net/http"
)

// Public Methods

func SpotifyGet(w http.ResponseWriter, r *http.Request) {
	// Creates a spotify model with it's secrets.
	spotify := new(models.Spotify)
	spotify.InitSecrets()

	// Gets result for RefreshAndAccessTokes request.
	spotify.Authorization.Response.Code = r.URL.Query().Get("code")
	// TODO: Check if there is a code value. (that means user is login.)
	// Result is in spotify.ResponseRefreshAndAccessTokens
	err := spotify.GetRefreshAndAccessTokensResponse()
	if err != nil {
		panic("Something wrong")
	}

	// Gets current song.
	artistName, songName, err := spotify.GetCurrentlyPlaying()
	if err != nil {
		panic("Something wrong")
	}

	// Redirects to lyrics page.
	path := "?artistName=" + artistName + "&songName=" + songName
	http.Redirect(w, r, "/lyric"+path, http.StatusFound)
}
