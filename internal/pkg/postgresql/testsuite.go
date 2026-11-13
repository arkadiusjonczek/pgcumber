package postgresql

import (
	"context"
	"fmt"

	"github.com/arkadiusjonczek/pgcumber/internal/pkg/cucumber"
	"github.com/cucumber/godog"
)

var _ cucumber.GodogTestsuite = (*Testsuite)(nil)

func NewTestsuite(
	autoconnect bool,
	host string,
	port string,
	username string,
	password string,
	database string,
) *Testsuite {
	return &Testsuite{
		autoconnect: autoconnect,
		host:        host,
		port:        port,
		username:    username,
		password:    password,
		database:    database,
	}
}

type Testsuite struct {
	autoconnect bool

	host     string
	port     string
	username string
	password string
	database string
}

func (suite *Testsuite) InitializeTestsuite(_ *godog.TestSuiteContext) {
	// nothing to do right now
}

func (suite *Testsuite) InitializeScenario(sc *godog.ScenarioContext) {
	scenario := &Scenario{
		host:     suite.host,
		port:     suite.port,
		username: suite.username,
		password: suite.password,
		database: suite.database,
	}

	sc.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		for _, tag := range sc.Tags {
			if tag.Name == "@autoconnect" {
				scenario.autoconnect = true
			}
		}

		if suite.autoconnect || scenario.autoconnect {
			// TODO: use specific function instead of cucumber function
			err := scenario.iConnectToTheServer()
			if err != nil {
				return ctx, fmt.Errorf("could not connect to postgresql server: %w", err)
			}
		}

		return ctx, nil
	})

	sc.Given(`^The host "([^"]*)"$`, scenario.theHost)
	sc.Given(`^The port "([^"]*)"$`, scenario.thePort)
	sc.Given(`^The database "([^"]*)"$`, scenario.theDatabase)

	sc.When(`^I use username "([^"]*)"$`, scenario.iUseUsername)
	sc.When(`^I use password "([^"]*)"$`, scenario.iUsePassword)

	sc.When(`^I use username from environment variable "([^"]*)"$`, scenario.iUseUsernameFromEnvironmentVariable)

	sc.When(`^I use password from environment variable "([^"]*)"$`, scenario.iUsePasswordFromEnvironmentVariable)

	sc.When(`^I use username from file "([^"]*)"$`, scenario.iUseUsernameFromFile)

	sc.When(`^I use password from file "([^"]*)"$`, scenario.iUsePasswordFromFile)

	sc.When(`^I connect to the server$`, scenario.iConnectToTheServer)

	sc.Then(`^the connection is successful$`, scenario.theConnectionIsSuccessful)

	sc.Then(`^I am a database superuser$`, scenario.iAmADatabaseSuperUser)
	sc.Then(`^I am not a database superuser$`, scenario.iAmNotADatabaseSuperUser)

	sc.Then(`^the schema "([^"]*)" exists$`, scenario.theSchemaExists)
	sc.Then(`^the schema "([^"]*)" does not exist$`, scenario.theSchemaNotExists)

	sc.Then(`^the current user has "([^"]*)" privilege on schema "([^"]*)"$`, scenario.theCurrentUserHasPrivilegeOnSchema)
	sc.Then(`^the user "([^"]*)" has "([^"]*)" privilege on schema "([^"]*)"$`, scenario.theUserHasPrivilegeOnSchema)

	sc.Then(`^the table "([^"]*)" exists in schema "([^"]*)"$`, scenario.theTableExistsInSchema)

	sc.Then(`^the current user has "([^"]*)" privilege on table "([^"]*)"$`, scenario.theCurrentUserHasPrivilegeOnTable)

	sc.After(func(ctx context.Context, _ *godog.Scenario, err error) (context.Context, error) {
		if scenario.conn != nil {
			psqlerr := scenario.conn.Close(ctx)
			if psqlerr != nil {
				return ctx, psqlerr
			}
			scenario.conn = nil
		}
		return ctx, err
	})
}
