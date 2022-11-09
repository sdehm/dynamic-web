package server

import (
	"embed"
	"net/http"

	"github.com/sdehm/go-morph/templates"
)

// public content
//go:embed public
var public embed.FS

type Server struct {
	templates *templates.Templates
}

func Start(templates *templates.Templates) *Server {
	server := &Server{
		templates: templates,
	}
	http.HandleFunc("/", indexHandler(server))
	http.Handle("/public/", http.FileServer(http.FS(public)))
	http.ListenAndServe(":8080", nil)
	return server
}

func indexHandler(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.templates.Pages["index"].Execute(w, nil)
	}
}
