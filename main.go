package main

import (
	"github.com/deffusion/boot-server/server"
	"log"
)

func main() {
	s, err := server.New()
	if err != nil {
		log.Fatal(err)
		return
	}

	s.Run()
}
