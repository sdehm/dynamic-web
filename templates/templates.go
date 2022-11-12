package templates

import (
	"html/template"
)

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
