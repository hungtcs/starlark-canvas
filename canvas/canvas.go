package canvas

import (
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

var Module = &starlarkstruct.Module{
	Name: "canvas",
	Members: starlark.StringDict{
		"Context": starlark.NewBuiltin("Context", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (_ starlark.Value, err error) {
			var width, height starlark.Int
			err = starlark.UnpackArgs(fn.Name(), args, kwargs, "width", &width, "height", &height)
			if err != nil {
				return nil, err
			}
			return NewContext(width, height)
		}),
	},
}
