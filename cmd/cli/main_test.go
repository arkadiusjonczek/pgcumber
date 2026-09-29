package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_run(t *testing.T) {
	os.Args = os.Args[0:1]
	err := run()

	require.Error(t, err)
	require.ErrorContains(t, err, "usage: ")
	require.ErrorContains(t, err, " <path to features directory>")
}
