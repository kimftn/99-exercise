package main

import (
	"log"

	"property-api/internal/app"
)

func main() {
	server := app.NewServer()

	log.Println("starting server on :3000")
	if err := server.Listen(":3000"); err != nil {
		log.Fatal(err)
	}
}
