package main

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"github.com/boratanrikulu/lirik.app/controllers"
	"github.com/boratanrikulu/lirik.app/controllers/api"
)

func main() {
	godotenv.Load()

	go cloneOrPullDatabase(os.Getenv("DATABASE_ADDRESS"))

	r := mux.NewRouter()
	r.HandleFunc("/", controllers.WelcomeGet).Methods("GET")
	r.HandleFunc("/logout", controllers.LogoutGet).Methods("GET")
	r.HandleFunc("/spotify", controllers.SpotifyGet).Methods("GET")
	r.PathPrefix("/assets").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))

	a := r.PathPrefix("/api").Subrouter()
	a.HandleFunc("/search", api.Search).Methods("POST")

	serve(r, "3000")
}

func serve(r *mux.Router, defaultPort string) {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort // Default port if not specified
	}

	log.Println("Server started at :" + port + ".")
	log.Fatalln(http.ListenAndServe(":"+port, r))
}

func cloneOrPullDatabase(databaseURL string) {
	if databaseURL == "" {
		log.Println("DatabaseURL is not set. Caching will not be working.")
		return
	}

	var errOutput bytes.Buffer
	if folderExists("./database/.git") {
		cmd := exec.Command("git", "pull", "origin", "master")
		cmd.Dir = "./database"
		cmd.Stderr = &errOutput

		err := cmd.Run()
		if err != nil {
			log.Println(err)
			return
		}

		return
	}

	cmd := exec.Command("git", "clone", databaseURL, "./database")
	cmd.Stderr = &errOutput

	err := cmd.Run()
	if err != nil {
		log.Println(err)
		return
	}

	ctx := context.Background()
	go syncDatabase(ctx)
}

func syncDatabase(ctx context.Context) {
	for {
		time.Sleep(5 * time.Minute)
		if !folderExists("./database/.git") {
			log.Println("There is not .git folder")
			break
		}

		flag := true
		for _, command := range [][]string{
			[]string{"git", "config", "user.email", "bora@heroku.com"},
			[]string{"git", "config", "user.name", "HEROKU"},
			[]string{"git", "pull", "origin", "master"},
			[]string{"git", "add", "."},
			[]string{"git", "commit", "-m", "Add new lyrics"},
			[]string{"git", "push", "origin", "master"},
		} {
			errOutput := bytes.Buffer{}
			cmd := exec.Command(command[0], command[1:]...)
			cmd.Dir = "./database"
			cmd.Stderr = &errOutput
			err := cmd.Run()
			if msg := errOutput.String(); err != nil && msg != "" {
				log.Println("[Database]", msg)
				flag = false
				break
			}
		}

		if flag {
			log.Println("[Database] is synced.")
		}
	}
}

func folderExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}
