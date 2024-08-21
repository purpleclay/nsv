package nsv

import (
	"context"
	"io"
	"os"
	"strings"

	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

type DevNull struct{}

func (DevNull) Read(_ []byte) (int, error) {
	return 0, io.EOF
}

func (DevNull) Write(p []byte) (int, error) {
	return len(p), nil
}

func (DevNull) Close() error {
	return nil
}

func exec(cmd string, env []string) error {
	p, err := syntax.NewParser().Parse(strings.NewReader(cmd), "")
	if err != nil {
		return err
	}

	ienv := os.Environ()
	ienv = append(ienv, env...)

	r, err := interp.New(
		interp.Params("-e"),
		interp.StdIO(os.Stdin, os.Stderr, os.Stderr),
		interp.OpenHandler(openHandler),
		interp.Env(expand.ListEnviron(ienv...)),
	)
	if err != nil {
		return err
	}

	return r.Run(context.Background(), p)
}

func openHandler(ctx context.Context, path string, flag int, perm os.FileMode) (io.ReadWriteCloser, error) {
	if path == "/dev/null" {
		return DevNull{}, nil
	}

	return interp.DefaultOpenHandler()(ctx, path, flag, perm)
}
