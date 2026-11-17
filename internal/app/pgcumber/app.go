package pgcumber

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/arkadiusjonczek/pgcumber/internal/pkg/postgresql"
	"github.com/cucumber/godog"
)

type AppOption func(*App)

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

func WithAutoconnect() AppOption {
	return func(app *App) {
		app.autoconnect = true
	}
}

func WithColor(color bool) AppOption {
	return func(app *App) {
		app.color = color
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

	color bool
}

func (app *App) Run() int {
	postgresqlTestsuite := postgresql.NewTestsuite(
		app.autoconnect,
		app.host,
		app.port,
		app.username,
		app.password,
		app.database,
	)

	opts := godog.Options{
		Format:   "pretty",
		Output:   &cleanPrettyWriter{out: os.Stdout},
		Paths:    []string{app.featuresPath},
		NoColors: true,
	}

	testsuite := godog.TestSuite{
		Name:                 "pgcumber",
		TestSuiteInitializer: postgresqlTestsuite.InitializeTestsuite,
		ScenarioInitializer:  postgresqlTestsuite.InitializeScenario,
		Options:              &opts,
	}

	status := testsuite.Run()

	return status
}

type cleanPrettyWriter struct {
	out io.Writer
}

// Regex findet ANSI-Codes + "# pfad/zur/datei.feature:12" am Zeilenende
var locationRegex = regexp.MustCompile(`(\x1b\[[0-9;]*m)?\s+#\s+[^\s]+:\d+(\x1b\[[0-9;]*m)?`)

func (w *cleanPrettyWriter) Write(p []byte) (int, error) {
	cleaned := locationRegex.ReplaceAllString(string(p), "")
	// Optional: Falls auch leere Zeilen oder bestimmte Debug-Prefixes weg sollen:
	cleaned = strings.ReplaceAll(cleaned, "  # ", "")
	
	_, err := w.out.Write([]byte(cleaned))
	// Wichtig: Die ursprüngliche Byte-Länge zurückgeben, sonst meldet Godog einen Short-Write-Fehler
	return len(p), err
}
