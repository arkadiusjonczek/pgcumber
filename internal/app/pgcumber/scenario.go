package pgcumber

import (
	"github.com/cucumber/godog"
)

type Scenario interface {
	InitializeScenario(sc *godog.ScenarioContext)
}
