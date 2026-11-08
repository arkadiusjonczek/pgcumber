package main

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_run_ReturnsError(t *testing.T) {
	os.Args = os.Args[0:1]
	err := run()

	require.Error(t, err)
	require.ErrorContains(t, err, "usage: ")
	require.ErrorContains(t, err, " <path to features file or directory>")
}

func Test_run_ReturnsVersion(t *testing.T) {
	oldStdout := os.Stdout
	defer func() { os.Stdout = oldStdout }()

	r, w, _ := os.Pipe()
	os.Stdout = w

	// simulate pgcumber --version
	os.Args = []string{
		os.Args[0],
		"--version",
	}
	err := run()

	w.Close()

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)

	consoleOutput := buf.String()

	require.NoError(t, err)
	require.Equal(t, "pgcumber dev (Commit: none, Date: unknown)\n", consoleOutput)
}
