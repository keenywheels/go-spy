package main

import (
	"log"

	"github.com/keenywheels/go-spy/internal/scheduler"
)

func main() {
	if err := scheduler.New().Run(); err != nil {
		log.Fatal(err)
	}
}
