package main

import (
	"log"

	"github.com/sdehm/go-morph/server"
	"github.com/sdehm/go-morph/templates"
)

func main() {
	logger := log.New(log.Writer(), "server: ", log.Flags())
	log.Fatal(server.Start(templates.New(), logger))
}
