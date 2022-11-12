package server

import (
	"embed"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/sdehm/go-morph/templates"
)

// public content
//go:embed public
var public embed.FS

type Server struct {
	templates *templates.Templates
	conn			net.Conn
}

func Start(templates *templates.Templates) *Server {
	server := &Server{
		templates: templates,
	}
	http.HandleFunc("/", indexHandler(server))
	http.Handle("/public/", http.FileServer(http.FS(public)))
	
	http.Handle("/ws", wsHandler(server))
	
	go clockTick(server)

	http.ListenAndServe(":8080", nil)


	return server
}

func indexHandler(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.templates.Pages["index"].Execute(w, nil)
	}
}

type morphData struct {
	Id string		`json:"id"`
	Html string `json:"html"`
}

// sends a message to the client every second with the current time updated
func clockTick(s *Server) {
	fmt.Println("starting clock tick")
	for tick := range time.Tick(time.Second) {
		t := tick.Format(time.RFC3339)
		fmt.Println(t)
		s.send(morphData{
			Id: "clock",
			Html: fmt.Sprintf("<p id=\"clock\" class=\"text-base text-gray-500\">%s</p>", t),
		})
	}
}

// Serialize the data to JSON and send it to the client
func (s *Server) send(m morphData) {
	if s.conn == nil {
		return
	}
	data, err := json.Marshal(m)
	if err != nil {
		fmt.Println(err)
	}
	
	err = wsutil.WriteServerText(s.conn, data)
	if err != nil {
		fmt.Println(err)
	}
}

// currently just echos back the message and prints it to the console
func wsHandler(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		s.conn = conn
		if err != nil {
			panic(err)
		}
	}
}