package studio

import (
	_ "embed"
	"html/template"
	"net/http"
)

//go:embed studio.html
var studioHTML string

func Handler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("studio").Parse(studioHTML)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}
