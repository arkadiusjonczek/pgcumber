package pgcumber

import (
	"fmt"
	"os"

	"github.com/cucumber/godog"
)

func NewApp(featuresPath string) (*App, error) {
	if featuresPath == "" {
		return nil, fmt.Errorf("no features path provided")
	}

	return &App{
		featuresPath: featuresPath,
	}, nil
}

type App struct {
	featuresPath string
}

func (app *App) InitializeTestsuite(_ *godog.TestSuiteContext) {
	// nothing to initialize right now
}

func (app *App) InitializeScenario(sc *godog.ScenarioContext) {
	postgresqlScenario := &PostgresqlScenario{}
	postgresqlScenario.InitializeScenario(sc)
}

func (app *App) Run() {
	opts := godog.Options{
		Format: "pretty",
		Paths:  []string{app.featuresPath},
	}

	status := godog.TestSuite{
		Name:                 "pgcumber",
		TestSuiteInitializer: app.InitializeTestsuite,
		ScenarioInitializer:  app.InitializeScenario,
		Options:              &opts,
	}.Run()

	if status != 0 {
		os.Exit(status)
	}
}
