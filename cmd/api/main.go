package main

import (
	"log"

	"github.com/khuchuz/go-clean-architecture-sql/server"
)

func main() {

	app := server.NewApp()

	if err := app.Run("8000"); err != nil {
		log.Fatalf("%s", err.Error())
	}
}
