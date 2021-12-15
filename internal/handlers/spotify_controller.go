package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"

	"github.com/boratanrikulu/lirik.app/internal/handlers/helpers"
	"github.com/boratanrikulu/lirik.app/internal/models"
	"github.com/boratanrikulu/lirik.app/pkg/lyrics"
	"github.com/boratanrikulu/lirik.app/pkg/meta"
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
		code := r.URL.Query().Get("code")
		sAuthed, err = spotify.S.Login(state, code)
		if err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		if stateCookie == nil || stateCookie.Value != state {
			log.Print("[ERROR] User's state and request state are not same.")
			errorMessages := []string{
				"Your state cookie and the response are not same.",
				"You might be under attack.",
			}

			helpers.ErrorPage(errorMessages, w)
			return
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
	var song spotify.Song
	if r.URL.Path == "/search" {
		song, err = sAuthed.Search(r.URL.Query().Get("q"))
		if err != nil {
			errorMessages := []string{
				"We can't find any related song for the query. 😔",
				"We are sorry about that.",
				"",
				"Try another song! 🙏",
			}

			helpers.ErrorPage(errorMessages, w)
			return
		}
	} else {
		song, err = sAuthed.GetCurrentlyPlaying()
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
	}

	showLyric(song, w, r)
}

type lyricsPageData struct {
	IsAvaible bool
	Artist    models.Artist
	Song      models.Song
}

func showLyric(song spotify.Song, w http.ResponseWriter, r *http.Request) {
	var wg sync.WaitGroup
	var m meta.Meta
	wg.Add(1)
	go func() {
		defer wg.Done()
		metaFinder := meta.NewFinder()
		_, m = metaFinder.GetMeta(song.ArtistName, song.AlbumName)
	}()

	finder := lyrics.NewFinder()
	found, l := finder.GetLyrics(song.ArtistName, song.Name)

	wg.Wait()
	genre := m.Genre
	if m.Style != "" {
		genre = fmt.Sprintf("%s (%s)", genre, m.Style)
	}

	pageData := lyricsPageData{
		IsAvaible: found,
		Artist: models.Artist{
			Name: song.ArtistName,
		},
		Song: models.Song{
			Name:             song.Name,
			Lyrics:           l,
			AlbumName:        song.AlbumName,
			AlbumGenre:       genre,
			AlbumReleaseDate: song.ReleaseDate,
			AlbumImage:       song.AlbumImage,
			AlbumTotalTracks: fmt.Sprint(song.TotalTracks),
		},
	}

	files := helpers.GetTemplateFiles("./views/songs.html")
	tmpl := template.Must(template.ParseFiles(files...))
	_ = tmpl.Execute(w, pageData)
}
