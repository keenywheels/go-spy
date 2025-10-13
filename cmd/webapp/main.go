package main

import (
	"log"

	"github.com/keenywheels/go-spy/internal/webapp"
)

func main() {
	if err := webapp.New().Run(); err != nil {
		log.Fatal(err)
	}
}
