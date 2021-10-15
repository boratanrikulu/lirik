package workers

import (
	"bytes"
	"context"
	"log"
	"os"
	"os/exec"
	"time"
)

func PrepareDatabase(databaseURL string) {
	if databaseURL == "" {
		log.Println("DatabaseURL ise not set. Caching will not be working.")
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
			[]string{"git", "pull", "-X", "theirs", "origin", "master"},
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
