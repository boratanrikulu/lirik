package main

import (
	"log"
	"net/http"
	"os"

	"github.com/boratanrikulu/s-lyrics/controllers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	loadEnv()
	r := mux.NewRouter()
	r.HandleFunc("/", controllers.WelcomeGet).Methods("GET")
	r.HandleFunc("/lyric", controllers.LyricGet).Methods("GET")
	r.HandleFunc("/spotify", controllers.SpotifyGet).Methods("GET")
	serve(r, "3000")
}

func serve(r *mux.Router, defaultPort string) {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort // Default port if not specified
	}
	http.ListenAndServe(":"+port, r)
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error occur while loading .env file")
	}
}
