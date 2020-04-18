package helpers

import (
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
