package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/arkadiusjonczek/pgcumber/internal/app/pgcumber"
)

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
		appOptions = append(appOptions, pgcumber.WithUser(username))
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
