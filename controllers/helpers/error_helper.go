package helpers

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

type PageData struct {
	ErrorMessages []string
}

func ErrorPage(errors []string, w http.ResponseWriter) {
	p := PageData{
		ErrorMessages: errors,
	}
	files := GetTemplateFiles("./views/error.html")
	tpml := template.Must(template.ParseFiles(files...))
	err := tpml.Execute(w, p)
	if err != nil {
		log.Printf("Error occur while rendering error page: %v", err)
	}
}

func WriteErrorToRes(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(struct {
		Error string
	}{
		Error: message,
	})
}
