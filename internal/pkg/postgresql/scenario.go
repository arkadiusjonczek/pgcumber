package postgresql

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

type Scenario struct {
	autoconnect bool

	conn *pgx.Conn

	host     string
	port     string
	username string
	password string
	database string
}

func (scenario *Scenario) theHost(host string) error {
	if host == "" {
		return fmt.Errorf("host is required")
	}

	scenario.host = host

	return nil
}

func (scenario *Scenario) thePort(port string) error {
	if port == "" {
		return fmt.Errorf("port is required")
	}

	scenario.port = port

	return nil
}

func (scenario *Scenario) theDatabase(database string) error {
	if database == "" {
		return fmt.Errorf("database is required")
	}

	scenario.database = database

	return nil
}

func (scenario *Scenario) iUseUsername(username string) error {
	if username == "" {
		return fmt.Errorf("username is required")
	}

	scenario.username = username

	return nil
}

func (scenario *Scenario) iUsePassword(password string) error {
	if password == "" {
		return fmt.Errorf("password is required")
	}

	scenario.password = password

	return nil
}

func (scenario *Scenario) iUseUsernameFromEnvironmentVariable(usernameEnvironmentVariable string) error {
	if usernameEnvironmentVariable == "" {
		return fmt.Errorf("username environment variable is required")
	}

	username := os.Getenv(usernameEnvironmentVariable)
	if username == "" {
		return fmt.Errorf("username is required")
	}

	scenario.username = username

	return nil
}

func (scenario *Scenario) iUsePasswordFromEnvironmentVariable(passwordEnvironmentVariable string) error {
	if passwordEnvironmentVariable == "" {
		return fmt.Errorf("password environment variable is required")
	}

	password := os.Getenv(passwordEnvironmentVariable)
	if password == "" {
		return fmt.Errorf("password is required")
	}

	scenario.password = password

	return nil
}

func (scenario *Scenario) iUseUsernameFromFile(filepath string) error {
	if filepath == "" {
		return fmt.Errorf("username file path is required")
	}

	data, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("error reading user from file '%s'", filepath)
	}

	scenario.username = string(data)

	return nil
}

func (scenario *Scenario) iUsePasswordFromFile(filepath string) error {
	if filepath == "" {
		return fmt.Errorf("password file path is required")
	}

	data, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("error reading password from file '%s'", filepath)
	}

	scenario.password = string(data)

	return nil
}

func (scenario *Scenario) iConnectToTheServer() error {
	if scenario.host == "" {
		return fmt.Errorf("host is required")
	} else if scenario.username == "" {
		return fmt.Errorf("username is required")
	} else if scenario.password == "" {
		return fmt.Errorf("password is required")
	} else if scenario.database == "" {
		return fmt.Errorf("database is required")
	}

	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", scenario.username, scenario.password, scenario.host, scenario.port, scenario.database)

	return scenario.iConnectoTo(url)
}

func (scenario *Scenario) iConnectoTo(url string) error {
	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		return fmt.Errorf("could not connect to database: %w", err)
	}

	scenario.conn = conn

	return nil
}

func (scenario *Scenario) theConnectionIsSuccessful() error {
	if scenario.conn == nil {
		return fmt.Errorf("the database connection is nil")
	}

	return nil
}

func (scenario *Scenario) iAmADatabaseSuperUser() error {
	if scenario.conn == nil {
		return fmt.Errorf("could not access database because connection is nil")
	}

	row := scenario.conn.QueryRow(context.Background(), "SELECT usesuper FROM pg_user WHERE usename = current_user;")

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

func (scenario *Scenario) iAmNotADatabaseSuperUser() error {
	if scenario.conn == nil {
		return fmt.Errorf("could not access database because connection is nil")
	}

	row := scenario.conn.QueryRow(context.Background(), "SELECT usesuper FROM pg_user WHERE usename = current_user;")

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

func (scenario *Scenario) theSchemaExists(schema string) error {
	if scenario.conn == nil {
		return fmt.Errorf("could not access database because connection is nil")
	}

	row := scenario.conn.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM information_schema.schemata WHERE schema_name = $1);", schema)

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

func (scenario *Scenario) theSchemaNotExists(schema string) error {
	if scenario.conn == nil {
		return fmt.Errorf("could not access database because connection is nil")
	}

	row := scenario.conn.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM information_schema.schemata WHERE schema_name = $1);", schema)

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

func (scenario *Scenario) theCurrentUserHasPrivilegeOnSchema(privilege string, schemaName string) error {
	if scenario.conn == nil {
		return fmt.Errorf("no active database connection")
	}

	var hasPrivilege bool
	query := `SELECT has_schema_privilege(current_user, $1, $2)`

	err := scenario.conn.QueryRow(context.Background(), query, schemaName, privilege).Scan(&hasPrivilege)
	if err != nil {
		return fmt.Errorf("failed to check %s privilege for current user on schema %q: %w", privilege, schemaName, err)
	}

	if !hasPrivilege {
		return fmt.Errorf("expected current user to have %s privilege on schema %q, but they do not", privilege, schemaName)
	}

	return nil
}

func (scenario *Scenario) theUserHasPrivilegeOnSchema(user string, privilege string, schemaName string) error {
	if scenario.conn == nil {
		return fmt.Errorf("no active database connection")
	}

	var hasPrivilege bool
	query := `SELECT has_schema_privilege($1, $2, $3)`

	err := scenario.conn.QueryRow(context.Background(), query, user, schemaName, privilege).Scan(&hasPrivilege)
	if err != nil {
		return fmt.Errorf("failed to check %s privilege for user %s on schema %q: %w", privilege, user, schemaName, err)
	}

	if !hasPrivilege {
		return fmt.Errorf("expected user %s to have %s privilege on schema %q, but they do not", user, privilege, schemaName)
	}

	return nil
}

func (scenario *Scenario) theTableExistsInSchema(tableName string, schemaName string) error {
	if scenario.conn == nil {
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

	err := scenario.conn.QueryRow(context.Background(), query, schemaName, tableName).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if table %q.%q exists: %w", schemaName, tableName, err)
	}

	if !exists {
		return fmt.Errorf("expected table %q to exist in schema %q, but it does not", tableName, schemaName)
	}

	return nil
}

func (scenario *Scenario) theCurrentUserHasPrivilegeOnTable(privilege string, tableName string) error {
	if scenario.conn == nil {
		return fmt.Errorf("no active database connection")
	}

	var hasPrivilege bool
	// tableName kann hier direkt als "public.users" übergeben werden
	query := `SELECT has_table_privilege(current_user, $1, $2)`

	err := scenario.conn.QueryRow(context.Background(), query, tableName, privilege).Scan(&hasPrivilege)
	if err != nil {
		return fmt.Errorf("failed to check %s privilege on %q: %w", privilege, tableName, err)
	}
	if !hasPrivilege {
		return fmt.Errorf("expected current user to have %s on %q", privilege, tableName)
	}

	return nil
}
