package main

import (
	"app/main/internal/app"
	"log"
)

func main() {

	a := app.New()

	if err := a.Init(); err != nil {
		log.Fatal(err)
	}

	if err := a.Run(); err != nil {
		log.Fatal(err)
	}
}
