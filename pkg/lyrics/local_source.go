package lyrics

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/boratanrikulu/lirik.app/pkg/lyrics/constants"
)

type localSource struct{}

func newLocalSource() *localSource {
	return &localSource{}
}

func (f *localSource) GetLyrics(artistName string, songName string) (found bool, lyrics Lyrics) {
	fileName := getFileName(artistName, songName)
	file, err := os.Open(fileName)
	if err != nil {
		return
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&lyrics)
	if err != nil {
		return
	}

	lyrics.Translates = f.getTranslations(fileName)
	return true, lyrics
}

func (f *localSource) getTranslations(fileName string) (translates []Translate) {
	for _, translate := range constants.AllowedTranslationLanguages {
		fName := fileName + "_" + translate

		f, err := os.Open(fName)
		if err != nil {
			continue
		}
		defer f.Close()

		t := Translate{}
		err = json.NewDecoder(f).Decode(&t)
		if err != nil {
			continue
		}

		translates = append(translates, t)
	}

	return translates
}

func getFileName(artist string, songName string) string {
	fileName := artist + "-" + songName + ".json"
	fileName = strings.ReplaceAll(fileName, "/", "_")
	fileName = strings.ReplaceAll(fileName, "\\", "_")
	fileName = "./database/lyrics/" + fileName

	return fileName
}

func contains(array []string, value string) bool {
	for _, v := range array {
		if v == value {
			return true
		}
	}
	return false
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func saveToFile(artistName, songName string, l Lyrics) {
	fileName := getFileName(artistName, songName)
	f, err := os.Create(fileName)
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()

	b, _ := json.Marshal(l)
	_, err = f.Write(b)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("[CREATED] %s\n", fileName)
	if len(l.Translates) != 0 {
		go saveTranslationsToFiles(fileName, l)
	}
}

func saveTranslationsToFiles(fileName string, l Lyrics) {
	for _, translate := range l.Translates {
		fName := fileName + "_" + translate.Language
		if fileExists(fName) {
			continue
		}

		f, err := os.Create(fName)
		if err != nil {
			log.Println(err)
			return
		}
		defer f.Close()

		b, _ := json.Marshal(translate)
		_, err = f.Write(b)
		if err != nil {
			log.Println(err)
			return
		}

		log.Printf("[CREATED %s translation] %s\n", translate.Language, fName)
	}
}
