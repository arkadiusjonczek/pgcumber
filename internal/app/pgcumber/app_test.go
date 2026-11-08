package pgcumber

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_NewApp_ReturnsError(t *testing.T) {
	app, err := NewApp("")

	require.Error(t, err)
	require.ErrorContains(t, err, "no features path provided")

	require.Nil(t, app)
}

func Test_NewApp(t *testing.T) {
	app, err := NewApp(filepath.Join("../../../", "features/"))

	require.NoError(t, err)
	require.NotNil(t, app)
	require.IsType(t, &App{}, app)
}
