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

	dockerCompose, dockerComposeErr := compose.NewDockerCompose(filepath.Join("../../", "docker-compose.yaml"))
	require.NoError(t, dockerComposeErr)

	dockerComposeUpErr := dockerCompose.
		WaitForService("postgres", wait.ForListeningPort("5432/tcp")).
		Up(ctx, compose.Wait(true))
	require.NoError(t, dockerComposeUpErr)

	defer func() {
		dockerComposeDownErr := dockerCompose.Down(ctx)
		require.NoError(t, dockerComposeDownErr)
	}()

	app, err := pgcumber.NewApp(filepath.Join("../../", "features/"))
	require.NoError(t, err)
	require.NotNil(t, app)

	status := app.Run()
	require.Equal(t, 0, status)
}
