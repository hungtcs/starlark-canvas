package examples

import (
	"testing"

	"github.com/hungtcs/starlark-canvas/canvas"

	"go.starlark.net/starlark"
	"go.starlark.net/syntax"
)

func TestExecMain(t *testing.T) {
	thread := &starlark.Thread{}
	predeclared := starlark.StringDict{
		"canvas": canvas.Module,
	}
	_, err := starlark.ExecFileOptions(&syntax.FileOptions{TopLevelControl: true}, thread, "main.star", nil, predeclared)
	if err != nil {
		panic(err)
	}
}
