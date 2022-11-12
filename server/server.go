package server

import (
	"embed"
	"fmt"
	"net/http"
	"time"

	"github.com/gobwas/ws"
	"github.com/sdehm/go-morph/templates"
)

// public content
//go:embed public
var public embed.FS

type Server struct {
	templates          *templates.Templates
	connections        []*connection
	connectionsPending chan *connection
}

func Start(templates *templates.Templates) *Server {
	server := &Server{
		templates:          templates,
		connectionsPending: make(chan *connection),
	}
	http.HandleFunc("/", indexHandler(server))
	http.Handle("/public/", http.FileServer(http.FS(public)))

	http.Handle("/ws", wsHandler(server))

	go clockTick(server)
	go server.startConnectionAdder()

	http.ListenAndServe(":8080", nil)

	return server
}

func indexHandler(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.templates.Pages["index"].Execute(w, nil)
	}
}

type morphData struct {
	Id   string `json:"id"`
	Html string `json:"html"`
}

// sends a message to the client every second with the current time updated
func clockTick(s *Server) {
	fmt.Println("starting clock tick")
	for tick := range time.Tick(time.Second) {
		t := tick.Format(time.RFC3339)
		fmt.Println(t)
		s.broadcast(morphData{
			Id:   "clock",
			Html: fmt.Sprintf("<p id=\"clock\" class=\"text-base text-gray-500\">%s</p>", t),
		})
	}
}

// currently just echos back the message and prints it to the console
func wsHandler(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		s.addConnection(newConnection(conn))
		if err != nil {
			panic(err)
		}
	}
}

func (s *Server) broadcast(m morphData) {
	for _, c := range s.connections {
		c.send(m)
	}
}

func (s *Server) startConnectionAdder() {
	for c := range s.connectionsPending {
		s.connections = append(s.connections, c)
	}
}

func (s *Server) addConnection(c *connection) {
	s.connectionsPending <- c
}
