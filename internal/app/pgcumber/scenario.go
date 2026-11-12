package pgcumber

import (
	"github.com/cucumber/godog"
)

type Testsuite interface {
	InitializeTestsuite(tsc *godog.TestSuiteContext)
}

type Scenario interface {
	InitializeScenario(sc *godog.ScenarioContext)
}
