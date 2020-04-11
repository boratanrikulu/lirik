package main

import (
	"net/http"
	"os"

	"github.com/boratanrikulu/s-lyrics/controllers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	r := mux.NewRouter()
	r.HandleFunc("/spotify", controllers.SpotifyGet).Methods("GET")
	r.HandleFunc("/", controllers.WelcomeGet).Methods("GET")
	r.HandleFunc("/lyric", controllers.LyricGet).Methods("GET")
	r.HandleFunc("/logout", controllers.LogoutGet).Methods("GET")
	serve(r, "3000")
}

func serve(r *mux.Router, defaultPort string) {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort // Default port if not specified
	}
	http.ListenAndServe(":"+port, r)
}
