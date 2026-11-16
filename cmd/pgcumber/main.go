package main

import (
	"fmt"
	"os"

	"github.com/arkadiusjonczek/pgcumber/internal/app/pgcumber"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var (
	envPGPASSFILEKey = "PGPASSFILE"
	envPGHOSTKey     = "PGHOST"
	envPGPORTKey     = "PGPORT"
	envPGUSERKey     = "PGUSER"
	envPGPASSWORDKey = "PGPASSWORD"
	envPGDATABASEKey = "PGDATABASE"
)

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}
}

func run() error {
	if len(os.Args) < 2 {
		fmt.Println(help())
		return fmt.Errorf("no arguments given")
	}

	appOptions := []pgcumber.AppOption{}
	featuresPath := ""

	switch os.Args[1] {
	case "--help":
		fmt.Println(help())
		return nil
	case "--version":
		fmt.Printf("pgcumber %s (Commit: %s, Date: %s)\n", version, commit, date)
		return nil
	case "--autoconnect":
		if len(os.Args) < 3 {
			fmt.Println(help())
			return fmt.Errorf("no path to features file or directory specified")
		}
		featuresPath = os.Args[2]
		appOptions = append(appOptions, pgcumber.WithAutoconnect())
	default:
		featuresPath = os.Args[1]
	}

	appOptions, pgPassfileErr := readPGPassFile(appOptions)
	if pgPassfileErr != nil {
		return fmt.Errorf("could not read pg passfiles: %v", pgPassfileErr)
	}

	appOptions, pgEnvErr := readPGEnv(appOptions)
	if pgEnvErr != nil {
		return fmt.Errorf("could not read environment variables: %v", pgEnvErr)
	}

	app, err := pgcumber.NewApp(featuresPath, appOptions...)
	if err != nil {
		return fmt.Errorf("could not create app: %w", err)
	}

	status := app.Run()
	if status != 0 {
		return fmt.Errorf("error running app: %d", status)
	}

	return nil
}

func help() string {
	var helpString string
	helpString += "Usage: pgcumber [option] <path to features file or directory>\n"
	helpString += "\n"
	helpString += "Use Gherkin to write plain-text, human-readable test cases for your PostgreSQL data\n"
	helpString += "\n"
	helpString += "Options:\n"
	helpString += "\t--autoconnect\tAuto connect to the PostgreSQL server (requires autoconfig)\n"
	helpString += "\t--version\tPrint version information and quit\n"
	return helpString
}
