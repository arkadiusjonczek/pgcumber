//go:build integration

package pgcumber

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/compose"
	"github.com/testcontainers/testcontainers-go/wait"
)

func Test_App_Run(t *testing.T) {
	ctx := context.Background()

	stack, err := compose.NewDockerCompose(filepath.Join("../../../", "docker-compose.yaml"))
	if err != nil {
		t.Fatal(err)
	}

	err = stack.
		WaitForService("postgres", wait.ForListeningPort("5432/tcp")).
		Up(ctx,
			compose.Wait(true),
		)
	require.NoError(t, err)

	defer stack.Down(ctx)

	// postgresContainer, err := testcontainers.Run(
	// 	ctx, "postgres:18.6-alpine3.24",
	// 	testcontainers.WithExposedPorts("5432/tcp"),
	// 	testcontainers.WithWaitStrategy(
	// 		wait.ForListeningPort("5432/tcp"),
	// 		// wait.ForLog("Ready to accept connections"),
	// 	),
	// 	testcontainers.WithEnv(map[string]string{
	// 		"POSTGRES_DB":       "pgcumber",
	// 		"POSTGRES_USER":     "admin",
	// 		"POSTGRES_PASSWORD": "admin",
	// 	}),
	// )
	// require.NoError(t, err)
	// testcontainers.CleanupContainer(t, postgresContainer)

	app, err := NewApp(filepath.Join("../../../", "features/"))
	require.NoError(t, err)
	require.NotNil(t, app)

	app.Run()
}
