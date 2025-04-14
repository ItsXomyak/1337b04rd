package http

import (
	"html/template"
	"net/http"
	"path/filepath"
)

var templates *template.Template

func InitTemplates() error {
	pattern := filepath.Join("web", "templates", "*.html")
	var err error
	templates, err = template.ParseGlob(pattern)
	if err != nil {
		return err
	}
	return nil
}

func RenderTemplate(w http.ResponseWriter, name string, data interface{}) error {
	return templates.ExecuteTemplate(w, name, data)
}