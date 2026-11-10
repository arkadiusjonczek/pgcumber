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

var (
	envPGHOSTKey     = "PGHOST"
	envPGPORTKey     = "PGPORT"
	envPGUSERKey     = "PGUSER"
	envPGPASSWORDKey = "PGPASSWORD"
	envPGDATABASEKey = "PGDATABASE"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func help() string {
	var helpString string
	helpString += "Usage: pgcumber [option] <path to features file or directory>\n"
	helpString += "\n"
	helpString += "Use Gherkin to write plain-text, human-readable test cases for your PostgreSQL data\n"
	helpString += "\n"
	helpString += "Options:\n"
	helpString += "\t--autoconnect\tConnect to the PostgreSQL database using DATABASE_URL\n"
	helpString += "\t--version\tPrint version information and quit\n"
	return helpString
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

	appOptions, err := readEnv(appOptions)
	if err != nil {
		return fmt.Errorf("could not read environment variables: %v", err)
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

func readEnv(appOptions []pgcumber.AppOption) ([]pgcumber.AppOption, error) {
	if appOptions == nil {
		return nil, fmt.Errorf("appOptions is nil")
	}

	pghost := os.Getenv(envPGHOSTKey)
	if len(pghost) > 0 {
		appOptions = append(appOptions, pgcumber.WithHost(pghost))
	}

	pgport := os.Getenv(envPGPORTKey)
	if len(pgport) > 0 {
		appOptions = append(appOptions, pgcumber.WithPort(pgport))
	}

	pguser := os.Getenv(envPGUSERKey)
	if len(pguser) > 0 {
		appOptions = append(appOptions, pgcumber.WithUser(pguser))
	}

	pgpassword := os.Getenv(envPGPASSWORDKey)
	if len(pgpassword) > 0 {
		appOptions = append(appOptions, pgcumber.WithPassword(pgpassword))
	}

	pgdatabase := os.Getenv(envPGDATABASEKey)
	if len(pgdatabase) > 0 {
		appOptions = append(appOptions, pgcumber.WithDatabase(pgdatabase))
	}

	return appOptions, nil
}
