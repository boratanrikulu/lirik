package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"github.com/boratanrikulu/lirik.app/internal/handlers"
	"github.com/boratanrikulu/lirik.app/internal/handlers/api"
	"github.com/boratanrikulu/lirik.app/internal/workers"
	"github.com/boratanrikulu/lirik.app/pkg/spotify"
)

func main() {
	godotenv.Load()

	go workers.PrepareDatabase(os.Getenv("DATABASE_ADDRESS"))
	r := mux.NewRouter()
	r.HandleFunc("/", handlers.WelcomeGet).Methods("GET")
	r.HandleFunc("/logout", handlers.LogoutGet).Methods("GET")
	r.HandleFunc("/spotify", handlers.SpotifyGet).Methods("GET")
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./assets/"))))
	a := r.PathPrefix("/api").Subrouter()
	a.HandleFunc("/search", api.Search).Methods("POST")
	spotify.S = spotify.NewSpotify(
		os.Getenv("CLIENT_ID"),
		os.Getenv("CLIENT_SECRET"),
		os.Getenv("REDIRECT_URI"),
	)

	serve(r, "3000")
}

func serve(r *mux.Router, defaultPort string) {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	log.Println("Server started at :" + port + ".")
	log.Fatalln(http.ListenAndServe(":"+port, r))
}
