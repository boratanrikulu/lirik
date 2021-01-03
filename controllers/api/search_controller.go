package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/boratanrikulu/lirik.app/controllers/helpers"
	"github.com/boratanrikulu/lirik.app/models"
)

func Search(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	err := auth(r)
	if err != nil {
		helpers.WriteErrorToRes(w, http.StatusUnauthorized, err.Error())
		return
	}

	artistName := strings.TrimSpace(r.URL.Query().Get("artistName"))
	songName := strings.TrimSpace(r.URL.Query().Get("songName"))
	err = validate(artistName, songName)
	if err != nil {
		helpers.WriteErrorToRes(w, http.StatusBadRequest, err.Error())
		return
	}

	l := new(models.Lyric)
	l.GetLyric(artistName, songName)

	if !l.IsAvaible {
		log.Printf("[API] [NOT FOUND] \"%v by %v\"", songName, artistName)
		helpers.WriteErrorToRes(w, http.StatusOK, "We couldn't find any lyrics.")
		return
	}

	log.Printf("[API] [FOUND] \"%v by %v\"", songName, artistName)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(struct {
		Lines []string
	}{
		Lines: l.Lines,
	})
}

func auth(r *http.Request) error {
	key := os.Getenv("SEARCH_API_KEY")
	if key == "" {
		return nil
	}

	if r.Header.Get("api-key") == key {
		return nil
	}

	return errors.New("Auth is not valid to use Search API.")
}

func validate(artistName string, songName string) error {
	if artistName == "" || songName == "" {
		return errors.New("Artist name or song name can't be empty.")
	}

	return nil
}
