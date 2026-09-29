package pgcumber

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_NewApp(t *testing.T) {
	app, err := NewApp("")

	require.Error(t, err)
	require.ErrorContains(t, err, "no features path provided")

	require.Nil(t, app)
}
