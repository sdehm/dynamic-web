package server

import (
	"embed"
	"fmt"
	"net/http"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
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
	
	http.Handle("/ws", wsHandler(server))
	
	http.ListenAndServe(":8080", nil)


	return server
}

func indexHandler(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.templates.Pages["index"].Execute(w, nil)
	}
}

// currently just echos back the message and prints it to the console
func wsHandler(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			panic(err)
		}
		go func() {
			defer conn.Close()
			for {
				msg, err := wsutil.ReadClientText(conn)
				if err != nil {
					panic(err)
				}
				fmt.Println(string(msg))
				err = wsutil.WriteServerText(conn, msg)
				if err != nil {
					panic(err)
				}
			}
		}()
	}
}