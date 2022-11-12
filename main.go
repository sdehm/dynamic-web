package main

import (
	"log"

	"github.com/sdehm/go-morph/server"
	"github.com/sdehm/go-morph/templates"
)

func main() {
	log.Fatal(server.Start(templates.New()))
}
