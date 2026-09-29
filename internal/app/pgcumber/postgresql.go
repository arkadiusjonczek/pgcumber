package pgcumber

import (
	"context"
	"fmt"
	"os"

	"github.com/cucumber/godog"
	"github.com/jackc/pgx/v5"
)

var _ Scenario = (*Postgresql)(nil)

type Postgresql struct {
	conn *pgx.Conn

	host     string
	port     string
	username string
	password string
	database string
}

func (psql *Postgresql) InitializeScenario(sc *godog.ScenarioContext) {
	sc.Given(`^The host "([^"]*)"$`, psql.theHost)
	sc.Given(`^The port "([^"]*)"$`, psql.thePort)
	sc.Given(`^The database "([^"]*)"$`, psql.theDatabase)

	sc.When(`^I use username "([^"]*)"$`, psql.iUseUsername)
	sc.When(`^I use password "([^"]*)"$`, psql.iUsePassword)
	sc.When(`^I use password from environment variable "([^"]*)"$`, psql.iUsePasswordFromEnvironmentVariable)
	sc.When(`^I connect to the server$`, psql.iConnectToTheServer)

	sc.Then(`^the connection is successful$`, psql.theConnectionIsSuccessful)
	sc.Then(`^I am a database superuser$`, psql.iAmADatabaseSuperUser)
	sc.Then(`^I am not a database superuser$`, psql.iAmNotADatabaseSuperUser)

	sc.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		if psql.conn != nil {
			psqlerr := psql.conn.Close(ctx)
			if psqlerr != nil {
				return ctx, psqlerr
			}
		}
		return ctx, err
	})
}

func (psql *Postgresql) theHost(host string) error {
	if host == "" {
		return fmt.Errorf("host is required")
	}

	psql.host = host

	return nil
}

func (psql *Postgresql) thePort(port string) error {
	if port == "" {
		return fmt.Errorf("port is required")
	}

	psql.port = port

	return nil
}

func (psql *Postgresql) theDatabase(database string) error {
	if database == "" {
		return fmt.Errorf("database is required")
	}

	psql.database = database

	return nil
}

func (psql *Postgresql) iUseUsername(username string) error {
	if username == "" {
		return fmt.Errorf("username is required")
	}

	psql.username = username

	return nil
}

func (psql *Postgresql) iUsePassword(password string) error {
	if password == "" {
		return fmt.Errorf("password is required")
	}

	psql.password = password

	return nil
}

func (psql *Postgresql) iUsePasswordFromEnvironmentVariable(passwordEnvironmentVariable string) error {
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

func (psql *Postgresql) iConnectToTheServer() error {
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

func (psql *Postgresql) iConnectoTo(url string) error {
	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		return fmt.Errorf("could not connect to database: %w", err)
	}

	psql.conn = conn

	return nil
}

func (psql *Postgresql) theConnectionIsSuccessful() error {
	if psql.conn == nil {
		return fmt.Errorf("the database connection is nil")
	}

	return nil
}

func (psql *Postgresql) iAmADatabaseSuperUser() error {
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

func (psql *Postgresql) iAmNotADatabaseSuperUser() error {
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
