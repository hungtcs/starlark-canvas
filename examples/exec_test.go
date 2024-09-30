package examples

import (
	"testing"

	"github.com/hungtcs/starlark-canvas/canvas"

	"go.starlark.net/starlark"
	"go.starlark.net/syntax"
)

func TestExecMain(t *testing.T) {
	if _, err := runFile("./main.star"); err != nil {
		t.Fatal(err)
	}
}

func Test02(t *testing.T) {
	if _, err := runFile("./02.star"); err != nil {
		t.Fatal(err)
	}
}

func runFile(name string) (starlark.StringDict, error) {
	thread := &starlark.Thread{}
	predeclared := starlark.StringDict{
		"canvas": canvas.Module,
	}
	return starlark.ExecFileOptions(&syntax.FileOptions{TopLevelControl: true}, thread, name, nil, predeclared)
}
