package controllers

import (
	"io/ioutil"
	"net/url"
	"fmt"
	"github.com/boratanrikulu/s-lyrics/models"
	"html/template"
	"net/http"
)

// Public Methods

func WrongGet(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./views/gowrong.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		panic("Something is wrong")
	}
}

func SpotifyGet(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Spotify - GET")
	// Creates a spotify model with it's secrets.
	spotify := new(models.Spotify)
	spotify.InitSecrets()

	accessToken, err := r.Cookie("AccessToken")
	refreshToken, err := r.Cookie("RefreshToken")
	if err != nil  {
		// Gets result for RefreshAndAccessTokes request.
		auth := &spotify.Authorization
		auth.Response.State = r.URL.Query().Get("state")
		state, _ := r.Cookie("State")
		if state.Value != auth.Response.State {
			// TODO: Add a error message that says
			// State code are not equal. Move!
			// http.Redirect(w, r, "/wrong", http.StatusSeeOther)
		}
		auth.Response.Code = r.URL.Query().Get("code")

		// Result is in spotify.ResponseRefreshAndAccessTokens
		err := spotify.GetRefreshAndAccessTokensResponse()
		if err != nil {
			// http.Redirect(w, r, "/wrong", http.StatusSeeOther)
		}
		setTokenCookies(w, spotify.RefreshAndAccessTokens)
	} else {
		tokensResponse := &spotify.RefreshAndAccessTokens.Response
		tokensResponse.AccessToken = accessToken.Value
		spotify.RefreshAndAccessTokens.Response.RefreshToken = refreshToken.Value
	}

	// Gets current song.
	artistName, songName, err := spotify.GetCurrentlyPlaying()
	if err != nil {
		http.Redirect(w, r, "/wrong", http.StatusSeeOther)
	}

	// Redirects to lyrics page.
	// u, _ := url.Parse("/lyric")
	// q, _ := url.ParseQuery(u.RawQuery)
	// q.Add("artistName", artistName)
	// q.Add("songName", songName)
	// u.RawQuery = q.Encode()
	// http.Redirect(w, r, fmt.Sprint(u), http.StatusSeeOtherd)

	u, _ := url.Parse("/lyric")
	q, _ := url.ParseQuery(u.RawQuery)
	q.Add("artistName", artistName)
	q.Add("songName", songName)
	u.RawQuery = q.Encode()

	fmt.Println("####################")
	fmt.Println(u)
	a := "http://localhost:3000" + fmt.Sprint(u)
	resp, _ := http.Get(a)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	w.Write(body)
}
