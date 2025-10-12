package main

import (
	"log"

	app "github.com/keenywheels/go-spy/internal/webapp"
)

func main() {
	if err := app.New().Run(); err != nil {
		log.Fatal(err)
	}
}
