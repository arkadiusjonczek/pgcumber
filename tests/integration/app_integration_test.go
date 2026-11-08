package integration

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/arkadiusjonczek/pgcumber/internal/app/pgcumber"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/compose"
	"github.com/testcontainers/testcontainers-go/wait"
)

func Test_App_Run(t *testing.T) {
	ctx := context.Background()

	dockerCompose, err := compose.NewDockerCompose(filepath.Join("../../", "docker-compose.yaml"))
	require.NoError(t, err)

	err = dockerCompose.
		WaitForService("postgres", wait.ForListeningPort("5432/tcp")).
		Up(ctx, compose.Wait(true))
	require.NoError(t, err)

	defer func() {
		err := dockerCompose.Down(ctx)
		require.NoError(t, err)
	}()

	app, err := pgcumber.NewApp(filepath.Join("../../", "features/"))
	require.NoError(t, err)
	require.NotNil(t, app)

	app.Run()
}
