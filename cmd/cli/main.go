package main

import (
	"fmt"
	"log"
	"os"

	"github.com/arkadiusjonczek/pgcumber/internal/app/pgcumber"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	if len(os.Args) < 2 {
		return fmt.Errorf("usage: %s <path to features directory>", os.Args[0])
	}

	featuresPath := os.Args[1]

	app, err := pgcumber.NewApp(featuresPath)
	if err != nil {
		return fmt.Errorf("could not create app: %w", err)
	}

	app.Run()

	return nil
}
