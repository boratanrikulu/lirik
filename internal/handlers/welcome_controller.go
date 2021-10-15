package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/boratanrikulu/lirik.app/internal/handlers/helpers"
	"github.com/boratanrikulu/lirik.app/pkg/spotify"
)

type WelcomePageData struct {
	SpotifyAuthLink string
}

func WelcomeGet(w http.ResponseWriter, r *http.Request) {
	accessTokenCookie, _ := r.Cookie("AccessToken")
	refreshTokenCookie, _ := r.Cookie("RefreshToken")
	if accessTokenCookie != nil || refreshTokenCookie != nil {
		http.Redirect(w, r, "/spotify", http.StatusSeeOther)
		return
	}

	files := helpers.GetTemplateFiles("./views/welcome.html")
	tmpl := template.Must(template.ParseFiles(files...))

	authLink, state, err := spotify.S.GetRequestAuthorizationLink()
	if err != nil {
		log.Println("Somethings are wrong. We are working on it.")
		helpers.ErrorPage([]string{"Somethings are wrong. We are working on it."}, w)
		return
	}
	helpers.SetStateCookie(state, w)
	_ = tmpl.Execute(w, WelcomePageData{
		SpotifyAuthLink: authLink,
	})
}
