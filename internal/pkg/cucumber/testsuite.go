package cucumber

import (
	"github.com/cucumber/godog"
)

type GodogTestsuite interface {
	InitializeTestsuite(tsc *godog.TestSuiteContext)
	InitializeScenario(sc *godog.ScenarioContext)
}
