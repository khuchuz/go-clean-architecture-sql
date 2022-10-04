package main

import (
	"log"

	"github.com/khuchuz/go-clean-architecture-sql/auth/app"
)

func main() {

	serv := app.NewApp()

	if err := serv.Run("8000"); err != nil {
		log.Fatalf("%s", err.Error())
	}
}
