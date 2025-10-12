package main

import (
	"log"

	"github.com/keenywheels/go-spy/internal/app"
)

func main() {
	if err := app.New().Run(); err != nil {
		log.Fatal(err)
	}
}
