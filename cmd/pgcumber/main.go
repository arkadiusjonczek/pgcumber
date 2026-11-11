package main

import (
	"fmt"
	"log"
	"os"
	"strings"

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

func readPGPassFile(appOptions []pgcumber.AppOption) ([]pgcumber.AppOption, error) {
	if appOptions == nil {
		return nil, fmt.Errorf("appOptions is nil")
	}

	pgpassfile := os.Getenv(envPGPASSFILEKey)
	if len(pgpassfile) == 0 {
		return appOptions, nil
	}

	// #nosec G703
	bytes, err := os.ReadFile(pgpassfile)
	if err != nil {
		return nil, fmt.Errorf("could not read %s specified in %s: %v", pgpassfile, envPGPASSFILEKey, err)
	}

	lines := strings.Split(string(bytes), "\n")
	firstLine := strings.TrimSpace(lines[0])

	configuration := strings.Split(firstLine, ":")

	if len(configuration) != 5 {
		return nil, fmt.Errorf("invalid configuration format of %s specified in %s", pgpassfile, envPGPASSFILEKey)
	}

	host := configuration[0]
	port := configuration[1]
	database := configuration[2]
	username := configuration[3]
	password := configuration[4]

	if len(host) > 0 {
		appOptions = append(appOptions, pgcumber.WithHost(host))
	}
	if len(port) > 0 {
		appOptions = append(appOptions, pgcumber.WithPort(port))
	}
	if len(database) > 0 {
		appOptions = append(appOptions, pgcumber.WithDatabase(database))
	}
	if len(username) > 0 {
		appOptions = append(appOptions, pgcumber.WithUser(database))
	}
	if len(password) > 0 {
		appOptions = append(appOptions, pgcumber.WithPassword(password))
	}

	return appOptions, nil
}

func readPGEnv(appOptions []pgcumber.AppOption) ([]pgcumber.AppOption, error) {
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
