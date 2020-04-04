package main

import (
  "os"
  "net/http"

  "github.com/boratanrikulu/s-lyrics/controllers"
  "github.com/gorilla/mux"
)


func main() {
  r := mux.NewRouter()
  r.HandleFunc("/", controllers.WelcomeGet).Methods("GET")
  r.HandleFunc("/a", controllers.WelcomePost).Methods("GET")
  r.HandleFunc("/spotify", controllers.SpotifyGet).Methods("GET")
  serve(r, "3000")
}

func serve(r *mux.Router, defaultPort string) {
  port := os.Getenv("PORT")
  if port == "" {
    port = defaultPort // Default port if not specified
  }
  http.ListenAndServe(":" + port, r)
}
