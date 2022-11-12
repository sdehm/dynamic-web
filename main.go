package main

import (
	"github.com/sdehm/go-morph/server"
	"github.com/sdehm/go-morph/templates"
)

func main() {
	server.Start(templates.New())
}
