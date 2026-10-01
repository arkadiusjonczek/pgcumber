package main

import (
	"fmt"
	"log"
	"os"

	"github.com/arkadiusjonczek/pgcumber/internal/app/pgcumber"
)

const (
	FeaturesPathEnvVar = "FEATURES_PATH"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	// goreleaser sets the features path
	// in the docker image this way
	if os.Getenv(FeaturesPathEnvVar) != "" {
		if len(os.Args) < 2 {
			os.Args = append(os.Args, os.Getenv(FeaturesPathEnvVar))
		} else {
			os.Args[1] = os.Getenv(FeaturesPathEnvVar)
		}
	}

	if len(os.Args) < 2 {
		return fmt.Errorf("usage: pgcumber <path to features file or directory>")
	}

	featuresPath := os.Args[1]

	app, err := pgcumber.NewApp(featuresPath)
	if err != nil {
		return fmt.Errorf("could not create app: %w", err)
	}

	app.Run()

	return nil
}
