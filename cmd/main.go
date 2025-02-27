package main

import (
	"log"

	"github.com/AlexNov03/UrlShortener/internal/app"
)

func main() {
	entryPoint := app.NewApiEntryPoint()
	err := entryPoint.Init()
	if err != nil {
		log.Fatalf("error while initializing app: %v", err)
	}

	err = entryPoint.Run()
	if err != nil {
		log.Fatalf("error while starting app: %v", err)
	}
}
