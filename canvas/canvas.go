package canvas

import (
	"fmt"

	"github.com/fogleman/gg"
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
		"radians": starlark.NewBuiltin("radians", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			var value starlark.Value
			if err := starlark.UnpackArgs(fn.Name(), args, kwargs, "degrees", &value); err != nil {
				return nil, err
			}
			var degrees float64
			switch value := value.(type) {
			case starlark.Int:
				degrees = float64(value.BigInt().Int64())
			case starlark.Float:
				degrees = float64(value)
			default:
				return nil, fmt.Errorf("degrees must be int or float, got %s", value.Type())
			}
			return starlark.Float(gg.Radians(degrees)), nil
		}),
		"degrees": starlark.NewBuiltin("radians", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			var value starlark.Value
			if err := starlark.UnpackArgs(fn.Name(), args, kwargs, "radians", &value); err != nil {
				return nil, err
			}
			var radians float64
			switch value := value.(type) {
			case starlark.Int:
				radians = float64(value.BigInt().Int64())
			case starlark.Float:
				radians = float64(value)
			default:
				return nil, fmt.Errorf("radians must be int or float, got %s", value.Type())
			}
			return starlark.Float(gg.Degrees(radians)), nil
		}),
	},
}
