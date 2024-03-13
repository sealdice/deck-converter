package main

import (
	"fmt"
	"io"
	"os"
)

func openInput(fn string) (io.ReadCloser, error) {
	f, err := os.Open(fn)
	if err != nil {
		pathErr := err.(*os.PathError) //nolint:errorlint
		return nil, fmt.Errorf("failed to open input file %q: %w", fn, pathErr.Err)
	}

	return f, nil
}

func openOutput(fn string) (io.WriteCloser, error) {
	opt := os.O_EXCL
	if flag.Overwrite {
		opt = os.O_TRUNC
	}

	f, err := os.OpenFile(fn, os.O_WRONLY|os.O_CREATE|opt, 0644)
	if err != nil {
		pathErr := err.(*os.PathError) //nolint:errorlint
		return nil, fmt.Errorf("failed to open output file %q: %w", fn, pathErr.Err)
	}

	return f, nil
}
