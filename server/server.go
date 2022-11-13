package server

import (
	"embed"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gobwas/ws"
	"github.com/sdehm/dynamic-web/templates"
)

// public content
//
//go:embed public
var public embed.FS

type Server struct {
	templates         *templates.Templates
	logger 					  *log.Logger
	connections       []*connection
	connectionUpdates chan func()
	lastId            int
}

func Start(templates *templates.Templates, logger *log.Logger) error {
	server := &Server{
		templates:         templates,
		logger:            logger,
		connectionUpdates: make(chan func()),
	}
	http.HandleFunc("/", server.indexHandler())
	http.Handle("/public/", http.FileServer(http.FS(public)))

	http.Handle("/ws", server.wsHandler())

	go server.clockTick()
	go server.startConnectionUpdates()

	return http.ListenAndServe(":8080", nil)
}

func (s *Server) indexHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.templates.Pages["index"].Execute(w, nil)
	}
}

func (s *Server) wsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		s.addConnection(conn)
		if err != nil {
			panic(err)
		}
	}
}

// sends a message to the client every second with the current time updated
func (s *Server) clockTick() {
	s.logger.Println("starting clock tick")
	for tick := range time.Tick(time.Second) {
		t := tick.Format(time.RFC3339)
		s.logger.Printf("sending time: %s", t)
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
		err := c.send(m)
		if err != nil {
			s.logger.Println(err)
			s.removeConnection(c)
		}
	}
}

func (s *Server) startConnectionUpdates() {
	for u := range s.connectionUpdates {
		u()
	}
}

func (s *Server) addConnection(c net.Conn) {
	s.connectionUpdates <- func() {
		s.lastId++
		s.connections = append(s.connections, newConnection(s.lastId, c))
		go s.updateConnectionCount()
	}
}

func (s *Server) removeConnection(c *connection) {
	s.connectionUpdates <- func() {
		for i, con := range s.connections {
			if con.id == c.id {
				s.connections = append(s.connections[:i], s.connections[i+1:]...)
				go s.updateConnectionCount()
				return
			}
		}
	}
}

func (s *Server) updateConnectionCount() {
	s.broadcast(morphData{
		Id:   "connections",
		Html: fmt.Sprintf("<p id=\"connections\" class=\"text-base text-gray-500\">%d connections </p>", len(s.connections)),
	})
}
