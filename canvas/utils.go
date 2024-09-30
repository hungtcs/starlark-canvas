package canvas

import (
	"fmt"
	"reflect"

	"github.com/fogleman/gg"
	"go.starlark.net/starlark"
)

func UnpackMethodArgs(name string, args starlark.Tuple, kwargs []starlark.Tuple, method reflect.Value) ([]reflect.Value, error) {
	methodType := method.Type()
	count := methodType.NumIn()

	values := make([]starlark.Value, count)
	inValuePointers := make([]any, count)
	for i := 0; i < count; i++ {
		inValuePointers[i] = &values[i]
	}
	if err := starlark.UnpackPositionalArgs(name, args, kwargs, count, inValuePointers...); err != nil {
		return nil, err
	}

	methodArgs := make([]reflect.Value, count)
	for i := 0; i < count; i++ {
		typ := methodType.In(i)
		value := values[i]

		switch typ.Kind() {
		case reflect.String:
			if str, ok := ValueToGoString(value); ok {
				methodArgs[i] = reflect.ValueOf(str)
			} else {
				return nil, fmt.Errorf("method %s args index %d expected string but got %s", name, i, value.Type())
			}
		case reflect.Int:
			var val int
			if v, ok := ValueToGoInt[int](value); ok {
				val = v
			} else {
				return nil, fmt.Errorf("method %s args index %d expected int but got %s", name, i, value.Type())
			}
			// 处理一些枚举类型
			switch typ.Name() {
			case "Align":
				methodArgs[i] = reflect.ValueOf(gg.Align(val))
			case "LineCap":
				methodArgs[i] = reflect.ValueOf(gg.LineCap(val))
			case "LineJoin":
				methodArgs[i] = reflect.ValueOf(gg.LineJoin(val))
			case "FillRule":
				methodArgs[i] = reflect.ValueOf(gg.FillRule(val))
			default:
				methodArgs[i] = reflect.ValueOf(val)
			}
		case reflect.Float64:
			if v, ok := ValueToGoFloat[float64](value); ok {
				methodArgs[i] = reflect.ValueOf(v)
			} else {
				return nil, fmt.Errorf("method %s args index %d expected float but got %s", name, i, value.Type())
			}
		default:
			return nil, fmt.Errorf("unsupported reflect method arg type %v", typ)
		}
	}

	return methodArgs, nil
}

// 多个参数则返回 starlark.Tuple
func PackMethodResults(method reflect.Value, results []reflect.Value) (starlark.Value, error) {
	methodType := method.Type()
	count := methodType.NumOut()
	if count < 1 {
		return starlark.None, nil
	}

	// 判断是否有错误参数
	lastOut := results[count-1]
	if lastOut.Type() == reflect.TypeFor[error]() {
		if !lastOut.IsNil() {
			return starlark.None, lastOut.Interface().(error)
		}
		results = results[:count-1]
	}

	// 剩下的参数
	var values = make([]starlark.Value, len(results))
	for idx, result := range results {
		if val, ok := GoValueToStarlarkValue(result.Interface()); ok {
			values[idx] = val
			continue
		}
		return nil, fmt.Errorf("unsupported method result value %v", result)
	}
	return starlark.Tuple(values), nil
}

func ValueToGoString(val starlark.Value) (_ string, ok bool) {
	switch val := val.(type) {
	case starlark.String:
		return string(val), true
	default:
		return "", false
	}
}

func ValueToGoInt[T ~int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64](val starlark.Value) (_ T, ok bool) {
	switch val := val.(type) {
	case starlark.Int:
		var v T
		if err := starlark.AsInt(val, &v); err != nil {
			return 0, false
		}
		return v, true
	default:
		return 0, false
	}
}

// 接受 int、float
func ValueToGoFloat[T ~float32 | float64](val starlark.Value) (_ T, ok bool) {
	switch val := val.(type) {
	case starlark.Int:
		return T(val.BigInt().Int64()), true
	case starlark.Float:
		return T(val), true
	default:
		return 0, false
	}
}

// support int*、float*、string、bool
func GoValueToStarlarkValue(val any) (_ starlark.Value, ok bool) {
	switch val := val.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return GoIntToStarlarkInt(val)
	case float32, float64:
		return GoFloatToStarlarkInt(val)
	case string:
		return starlark.String(val), true
	case bool:
		if val {
			return starlark.True, true
		} else {
			return starlark.False, true
		}
	default:
		return nil, false
	}
}

func GoIntToStarlarkInt(val any) (_ starlark.Int, ok bool) {
	var v int64
	switch val := val.(type) {
	case int:
		v = int64(val)
	case int8:
		v = int64(val)
	case int16:
		v = int64(val)
	case int32:
		v = int64(val)
	case int64:
		v = int64(val)
	case uint:
		v = int64(val)
	case uint8:
		v = int64(val)
	case uint16:
		v = int64(val)
	case uint32:
		v = int64(val)
	case uint64:
		v = int64(val)
	default:
		return starlark.MakeInt(0), false
	}
	return starlark.MakeInt64(v), true
}

func GoFloatToStarlarkInt(val any) (_ starlark.Float, ok bool) {
	switch val := val.(type) {
	case float32:
		return starlark.Float(float64(val)), true
	case float64:
		return starlark.Float(float64(val)), true
	default:
		return starlark.Float(0), false
	}
}
