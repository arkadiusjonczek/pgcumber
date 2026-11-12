package pgcumber

import (
	"fmt"

	"github.com/arkadiusjonczek/pgcumber/internal/pkg/postgresql"
	"github.com/cucumber/godog"
)

type AppOption func(*App)

func WithAutoconnect() AppOption {
	return func(app *App) {
		app.autoconnect = true
	}
}

func WithHost(host string) AppOption {
	return func(app *App) {
		app.host = host
	}
}

func WithPort(port string) AppOption {
	return func(app *App) {
		app.port = port
	}
}

func WithUser(username string) AppOption {
	return func(app *App) {
		app.username = username
	}
}

func WithPassword(password string) AppOption {
	return func(app *App) {
		app.password = password
	}
}

func WithDatabase(database string) AppOption {
	return func(app *App) {
		app.database = database
	}
}

func NewApp(featuresPath string, opts ...AppOption) (*App, error) {
	if featuresPath == "" {
		return nil, fmt.Errorf("no features path provided")
	}

	app := &App{
		featuresPath: featuresPath,
	}

	for _, opt := range opts {
		opt(app)
	}

	return app, nil
}

type App struct {
	featuresPath string

	host     string
	port     string
	username string
	password string
	database string

	autoconnect bool
}

func (app *App) Run() int {
	postgresqlTestsuite := postgresql.NewTestsuite(
		app.autoconnect,
		app.host,
		app.port,
		app.username,
		app.password,
		app.database,
	)

	opts := godog.Options{
		Format: "pretty",
		Paths:  []string{app.featuresPath},
	}

	testsuite := godog.TestSuite{
		Name:                 "pgcumber",
		TestSuiteInitializer: postgresqlTestsuite.InitializeTestsuite,
		ScenarioInitializer:  postgresqlTestsuite.InitializeScenario,
		Options:              &opts,
	}

	status := testsuite.Run()

	return status
}
