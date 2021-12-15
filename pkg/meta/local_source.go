package meta

import (
	"encoding/json"
	"log"
	"os"
	"strings"
)

type localSource struct{}

func newLocalSource() *localSource {
	return &localSource{}
}

func (f *localSource) GetMeta(artistName string, albumName string) (found bool, meta Meta) {
	fileName := getFileName(artistName, albumName)
	file, err := os.Open(fileName)
	if err != nil {
		return
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&meta)
	if err != nil {
		return
	}

	return true, meta
}

func getFileName(artist string, albumName string) string {
	fileName := artist + "-" + albumName + ".meta"
	fileName = strings.ReplaceAll(fileName, "/", "_")
	fileName = strings.ReplaceAll(fileName, "\\", "_")
	fileName = "./database/lyrics/" + fileName

	return fileName
}

func saveToFile(artistName, albumName string, m Meta) {
	fileName := getFileName(artistName, albumName)
	f, err := os.Create(fileName)
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()

	b, _ := json.Marshal(m)
	_, err = f.Write(b)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("[CREATED META] %s\n", fileName)
}
