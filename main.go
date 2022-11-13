package main

import (
	"log"

	"github.com/sdehm/dynamic-web/server"
	"github.com/sdehm/dynamic-web/templates"
)

func main() {
	logger := log.New(log.Writer(), "server: ", log.Flags())
	log.Fatal(server.Start(templates.New(), logger))
}
