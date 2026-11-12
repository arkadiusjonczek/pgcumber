package pgcumber

import (
	"context"
	"fmt"
	"os"

	"github.com/cucumber/godog"
	"github.com/jackc/pgx/v5"
)

var _ Scenario = (*PostgresqlScenario)(nil)

type PostgresqlScenario struct {
	autoconnect bool

	conn *pgx.Conn

	host     string
	port     string
	username string
	password string
	database string
}

func (psql *PostgresqlScenario) InitializeTestsuite(_ *godog.TestSuiteContext) {
	// nothing to do right now
}

func (psql *PostgresqlScenario) InitializeScenario(sc *godog.ScenarioContext) {
	sc.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		scenarioAutoconnect := false
		for _, tag := range sc.Tags {
			if tag.Name == "@autoconnect" {
				scenarioAutoconnect = true
			}
		}

		if psql.autoconnect || scenarioAutoconnect {
			// TODO: use specific function instead of cucumber function
			_ = psql.iConnectToTheServer()
		}

		return ctx, nil
	})

	sc.Given(`^The host "([^"]*)"$`, psql.theHost)
	sc.Given(`^The port "([^"]*)"$`, psql.thePort)
	sc.Given(`^The database "([^"]*)"$`, psql.theDatabase)

	sc.When(`^I use username "([^"]*)"$`, psql.iUseUsername)
	sc.When(`^I use password "([^"]*)"$`, psql.iUsePassword)

	sc.When(`^I use username from environment variable "([^"]*)"$`, psql.iUseUsernameFromEnvironmentVariable)

	sc.When(`^I use password from environment variable "([^"]*)"$`, psql.iUsePasswordFromEnvironmentVariable)

	sc.When(`^I use username from file "([^"]*)"$`, psql.iUseUsernameFromFile)

	sc.When(`^I use password from file "([^"]*)"$`, psql.iUsePasswordFromFile)

	sc.When(`^I connect to the server$`, psql.iConnectToTheServer)

	sc.Then(`^the connection is successful$`, psql.theConnectionIsSuccessful)

	sc.Then(`^I am a database superuser$`, psql.iAmADatabaseSuperUser)
	sc.Then(`^I am not a database superuser$`, psql.iAmNotADatabaseSuperUser)

	sc.Then(`^the schema "([^"]*)" exists$`, psql.theSchemaExists)
	sc.Then(`^the schema "([^"]*)" does not exist$`, psql.theSchemaNotExists)

	sc.Then(`^the user has "([^"]*)" privilege on schema "([^"]*)"$`, psql.theUserHasPrivilegeOnSchema)

	sc.Then(`^the table "([^"]*)" exists in schema "([^"]*)"$`, psql.theTableExistsInSchema)

	sc.After(func(ctx context.Context, _ *godog.Scenario, err error) (context.Context, error) {
		if psql.conn != nil {
			psqlerr := psql.conn.Close(ctx)
			if psqlerr != nil {
				return ctx, psqlerr
			}
			psql.conn = nil
		}
		return ctx, err
	})
}

func (psql *PostgresqlScenario) theHost(host string) error {
	if host == "" {
		return fmt.Errorf("host is required")
	}

	psql.host = host

	return nil
}

func (psql *PostgresqlScenario) thePort(port string) error {
	if port == "" {
		return fmt.Errorf("port is required")
	}

	psql.port = port

	return nil
}

func (psql *PostgresqlScenario) theDatabase(database string) error {
	if database == "" {
		return fmt.Errorf("database is required")
	}

	psql.database = database

	return nil
}

func (psql *PostgresqlScenario) iUseUsername(username string) error {
	if username == "" {
		return fmt.Errorf("username is required")
	}

	psql.username = username

	return nil
}

func (psql *PostgresqlScenario) iUsePassword(password string) error {
	if password == "" {
		return fmt.Errorf("password is required")
	}

	psql.password = password

	return nil
}

func (psql *PostgresqlScenario) iUseUsernameFromEnvironmentVariable(usernameEnvironmentVariable string) error {
	if usernameEnvironmentVariable == "" {
		return fmt.Errorf("username environment variable is required")
	}

	username := os.Getenv(usernameEnvironmentVariable)
	if username == "" {
		return fmt.Errorf("username is required")
	}

	psql.username = username

	return nil
}

func (psql *PostgresqlScenario) iUsePasswordFromEnvironmentVariable(passwordEnvironmentVariable string) error {
	if passwordEnvironmentVariable == "" {
		return fmt.Errorf("password environment variable is required")
	}

	password := os.Getenv(passwordEnvironmentVariable)
	if password == "" {
		return fmt.Errorf("password is required")
	}

	psql.password = password

	return nil
}

func (psql *PostgresqlScenario) iUseUsernameFromFile(filepath string) error {
	if filepath == "" {
		return fmt.Errorf("username file path is required")
	}

	data, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("error reading user from file '%s'", filepath)
	}

	psql.username = string(data)

	return nil
}

func (psql *PostgresqlScenario) iUsePasswordFromFile(filepath string) error {
	if filepath == "" {
		return fmt.Errorf("password file path is required")
	}

	data, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("error reading password from file '%s'", filepath)
	}

	psql.password = string(data)

	return nil
}

func (psql *PostgresqlScenario) iConnectToTheServer() error {
	if psql.host == "" {
		return fmt.Errorf("host is required")
	} else if psql.username == "" {
		return fmt.Errorf("username is required")
	} else if psql.password == "" {
		return fmt.Errorf("password is required")
	} else if psql.database == "" {
		return fmt.Errorf("database is required")
	}

	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", psql.username, psql.password, psql.host, psql.port, psql.database)

	return psql.iConnectoTo(url)
}

func (psql *PostgresqlScenario) iConnectoTo(url string) error {
	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		return fmt.Errorf("could not connect to database: %w", err)
	}

	psql.conn = conn

	return nil
}

func (psql *PostgresqlScenario) theConnectionIsSuccessful() error {
	if psql.conn == nil {
		return fmt.Errorf("the database connection is nil")
	}

	return nil
}

func (psql *PostgresqlScenario) iAmADatabaseSuperUser() error {
	if psql.conn == nil {
		return fmt.Errorf("could not access database because connection is nil")
	}

	row := psql.conn.QueryRow(context.Background(), "SELECT usesuper FROM pg_user WHERE usename = current_user;")

	var usesuper bool
	err := row.Scan(&usesuper)
	if err != nil {
		return fmt.Errorf("could not query database: %w", err)
	}

	if !usesuper {
		return fmt.Errorf("the current user is not a super user")
	}

	return nil
}

func (psql *PostgresqlScenario) iAmNotADatabaseSuperUser() error {
	if psql.conn == nil {
		return fmt.Errorf("could not access database because connection is nil")
	}

	row := psql.conn.QueryRow(context.Background(), "SELECT usesuper FROM pg_user WHERE usename = current_user;")

	var usesuper bool
	err := row.Scan(&usesuper)
	if err != nil {
		return fmt.Errorf("could not query database: %w", err)
	}

	if usesuper {
		return fmt.Errorf("the current user is a super user")
	}

	return nil
}

func (psql *PostgresqlScenario) theSchemaExists(schema string) error {
	if psql.conn == nil {
		return fmt.Errorf("could not access database because connection is nil")
	}

	row := psql.conn.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM information_schema.schemata WHERE schema_name = $1);", schema)

	var exists bool
	err := row.Scan(&exists)
	if err != nil {
		return fmt.Errorf("could not query database: %w", err)
	}

	if !exists {
		return fmt.Errorf("the schema '%s' does not exist", schema)
	}

	return nil
}

func (psql *PostgresqlScenario) theSchemaNotExists(schema string) error {
	if psql.conn == nil {
		return fmt.Errorf("could not access database because connection is nil")
	}

	row := psql.conn.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM information_schema.schemata WHERE schema_name = $1);", schema)

	var exists bool
	err := row.Scan(&exists)
	if err != nil {
		return fmt.Errorf("could not query database: %w", err)
	}

	if exists {
		return fmt.Errorf("the schema '%s' does exist", schema)
	}

	return nil
}

func (psql *PostgresqlScenario) theUserHasPrivilegeOnSchema(privilege string, schemaName string) error {
	if psql.conn == nil {
		return fmt.Errorf("no active database connection")
	}

	var hasPrivilege bool
	// Gültige Werte für privilege in PostgreSQL: 'USAGE', 'CREATE', 'USAGE, CREATE'
	query := `SELECT has_schema_privilege(current_user, $1, $2)`

	err := psql.conn.QueryRow(context.Background(), query, schemaName, privilege).Scan(&hasPrivilege)
	if err != nil {
		return fmt.Errorf("failed to check %s privilege for current user on schema %q: %w", privilege, schemaName, err)
	}

	if !hasPrivilege {
		return fmt.Errorf("expected current user to have %s privilege on schema %q, but they do not", privilege, schemaName)
	}

	return nil
}

func (psql *PostgresqlScenario) theTableExistsInSchema(tableName, schemaName string) error {
	if psql.conn == nil {
		return fmt.Errorf("no active database connection")
	}

	var exists bool
	query := `
		SELECT EXISTS (
			SELECT 1 
			FROM information_schema.tables 
			WHERE table_schema = $1 
			  AND table_name = $2
			  AND table_type = 'BASE TABLE'
		)`

	err := psql.conn.QueryRow(context.Background(), query, schemaName, tableName).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if table %q.%q exists: %w", schemaName, tableName, err)
	}

	if !exists {
		return fmt.Errorf("expected table %q to exist in schema %q, but it does not", tableName, schemaName)
	}

	return nil
}
