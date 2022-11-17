package templates

import (
	"embed"
	"html/template"
)

// templates
//
//go:embed *.go.html
var _ embed.FS

type Templates struct {
	Pages map[string]*template.Template
}

func New() *Templates {
	index, err := template.ParseFiles("templates/index.go.html")
	if err != nil {
		panic(err)
	}
	return &Templates{
		Pages: map[string]*template.Template{
			"index": index,
		},
	}
}
