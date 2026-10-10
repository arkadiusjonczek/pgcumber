package main

import (
	"fmt"
	"log"
	"os"

	"github.com/arkadiusjonczek/pgcumber/internal/app/pgcumber"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	if len(os.Args) < 2 {
		return fmt.Errorf("usage: pgcumber <path to features file or directory>")
	}

	if os.Args[1] == "--version" {
		fmt.Printf("pgcumber %s (Commit: %s, Date: %s)\n", version, commit, date)
		return nil
	}

	featuresPath := os.Args[1]

	app, err := pgcumber.NewApp(featuresPath)
	if err != nil {
		return fmt.Errorf("could not create app: %w", err)
	}

	app.Run()

	return nil
}
