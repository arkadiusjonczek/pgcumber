package pgcumber

import (
	"bytes"
	"io"
	"os"
	"regexp"
)

const (
	GodogPrettyLineDebugRegex = `[ \t]*#[^\r\n]*`
)

func NewCleanPrettyWriter() *cleanPrettyWriter {
	return &cleanPrettyWriter{
		regex: regexp.MustCompile(GodogPrettyLineDebugRegex),
		out:   os.Stdout,
	}
}

type cleanPrettyWriter struct {
	regex *regexp.Regexp
	out   io.Writer
	buf   bytes.Buffer
}

func (w *cleanPrettyWriter) Write(p []byte) (int, error) {
	w.buf.Write(p)

	for {
		line, err := w.buf.ReadBytes('\n')
		if err != nil {
			w.buf.Write(line)
			break
		}
		cleaned := w.regex.ReplaceAll(line, []byte(""))
		if _, writeErr := os.Stdout.Write(cleaned); writeErr != nil {
			return 0, writeErr
		}
	}

	return len(p), nil
}
