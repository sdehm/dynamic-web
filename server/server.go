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
	http.HandleFunc("/", server.indexHandler())
	http.Handle("/public/", http.FileServer(http.FS(public)))

	http.Handle("/ws", server.wsHandler())

	go server.clockTick()
	go server.startConnectionAdder()

	http.ListenAndServe(":8080", nil)

	return server
}

func (s *Server) indexHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.templates.Pages["index"].Execute(w, nil)
	}
}

func (s *Server) wsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		s.addConnection(newConnection(conn))
		if err != nil {
			panic(err)
		}
	}
}

// sends a message to the client every second with the current time updated
func (s *Server) clockTick() {
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

type morphData struct {
	Id   string `json:"id"`
	Html string `json:"html"`
}

func (s *Server) broadcast(m morphData) {
	for _, c := range s.connections {
		c.send(m)
	}
}

func (s *Server) startConnectionAdder() {
	for c := range s.connectionsPending {
		s.connections = append(s.connections, c)
		s.broadcast(morphData{
			Id:   "connections",
			Html: fmt.Sprintf("<p id=\"connections\" class=\"text-base text-gray-500\">%d connections </p>", len(s.connections)),
		})
	}
}

func (s *Server) addConnection(c *connection) {
	s.connectionsPending <- c
}
