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

func (a *App) InitializeScenario(sc *godog.ScenarioContext) {
	scenarios := []Scenario{
		&Postgresql{},
	}

	for _, scenario := range scenarios {
		scenario.InitializeScenario(sc)
	}
}

func (a *App) Run() {
	opts := godog.Options{
		Format: "pretty",
		Paths:  []string{a.featuresPath},
	}

	status := godog.TestSuite{
		Name:                "godog",
		ScenarioInitializer: a.InitializeScenario,
		Options:             &opts,
	}.Run()

	if status != 0 {
		os.Exit(status)
	}
}
