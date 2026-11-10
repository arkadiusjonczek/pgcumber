package pgcumber

import (
	"fmt"

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

	postgresqlScenario *PostgresqlScenario
}

func (app *App) initializeTestsuite(tsc *godog.TestSuiteContext) {
	app.postgresqlScenario.InitializeTestsuite(tsc)
}

func (app *App) initializeScenario(sc *godog.ScenarioContext) {
	app.postgresqlScenario.InitializeScenario(sc)
}

func (app *App) Run() int {
	app.postgresqlScenario = &PostgresqlScenario{
		autoconnect: app.autoconnect,
		host:        app.host,
		port:        app.port,
		username:    app.username,
		password:    app.password,
		database:    app.database,
	}

	opts := godog.Options{
		Format: "pretty",
		Paths:  []string{app.featuresPath},
	}

	status := godog.TestSuite{
		Name:                 "pgcumber",
		TestSuiteInitializer: app.initializeTestsuite,
		ScenarioInitializer:  app.initializeScenario,
		Options:              &opts,
	}.Run()

	return status
}
