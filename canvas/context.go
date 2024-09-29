package canvas

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/jpeg"
	"image/png"
	"reflect"

	"github.com/fogleman/gg"
	changecase "github.com/ku/go-change-case"
	"go.starlark.net/starlark"
)

// 此处列出来的方法对 starlark 不可见
var hiddenMethods = map[string]bool{"save_png": true, "load_font_face": true}

// 名称映射，默认使用 changecase.Pascal 将 starlark 名称转为 Go 名称，
// 但是有些名称无法适用这个规则，此处可以手动指定
var drawContextAttrNameMapping = map[string]string{
	"set_rgb": "SetRGB",
}

// 扩展 gg.Context 上面不存在的方法
var drawContextExpand = starlark.StringDict{
	"get_base64": starlark.NewBuiltin("get_base64", func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		d := b.Receiver().(*DrawContext)
		var (
			format   starlark.String = "png"
			_quality starlark.Int
		)
		if err := starlark.UnpackArgs(b.Name(), args, kwargs, "format??", &format, "quality??", &_quality); err != nil {
			return nil, err
		}
		var quality int
		if err := starlark.AsInt(_quality, &quality); err != nil {
			return nil, err
		}
		return d.GetBase64(string(format), quality)
	}),
	"get_data_uri": starlark.NewBuiltin("get_data_uri", func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		d := b.Receiver().(*DrawContext)
		var (
			format   starlark.String = "png"
			_quality starlark.Int
		)
		if err := starlark.UnpackArgs(b.Name(), args, kwargs, "format??", &format, "quality??", &_quality); err != nil {
			return nil, err
		}
		var quality int
		if err := starlark.AsInt(_quality, &quality); err != nil {
			return nil, err
		}
		return d.GetDataURI(string(format), quality)
	}),
	// "save_file": starlark.NewBuiltin("save_file", func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	// 	dc := b.Receiver().(*DrawContext).dc
	// 	buf := bytes.NewBuffer(nil)
	// 	if err := png.Encode(buf, dc.Image()); err != nil {
	// 		return nil, err
	// 	}
	// 	if err := os.WriteFile("out.png", buf.Bytes(), 0644); err != nil {
	// 		return nil, err
	// 	}
	// 	return starlark.None, nil
	// }),
}

type DrawContext struct {
	dc     *gg.Context
	frozen bool
}

// Attr implements starlark.HasAttrs.
func (d *DrawContext) Attr(name string) (starlark.Value, error) {
	// 首先查找内置的扩展
	if val, ok := drawContextExpand[name]; ok {
		if b, ok := val.(*starlark.Builtin); ok {
			return b.BindReceiver(d), nil
		}
		return val, nil
	}

	// 是否是隐藏的方法
	if hiddenMethods[name] {
		return nil, starlark.NoSuchAttrError(fmt.Sprintf("%s has no .%s field", d.Type(), name))
	}

	// 映射或转换名称
	var attrName string
	if value, ok := drawContextAttrNameMapping[name]; ok {
		attrName = value
	} else {
		attrName = changecase.Pascal(name)
	}

	// gg.Context 的反射对象
	rv := reflect.ValueOf(d.dc)
	rt := reflect.TypeOf(d.dc)

	// gg.Context 中没有公开属性，所以只处理方法
	method, ok := rt.MethodByName(attrName)
	if ok && method.IsExported() {
		method := rv.MethodByName(attrName)
		return starlark.NewBuiltin(name, func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			methodArgs, err := UnpackMethodArgs(name, args, kwargs, method)
			if err != nil {
				return nil, err
			}
			// 调用函数
			results := method.Call(methodArgs)
			return PackMethodResults(method, results)
		}), nil
	}

	// attr 不存在
	return nil, starlark.NoSuchAttrError(fmt.Sprintf("%s has no .%s field", d.Type(), name))
}

// AttrNames implements starlark.HasAttrs.
func (d *DrawContext) AttrNames() []string {
	names := make([]string, 0)
	// expands
	for name := range drawContextExpand {
		names = append(names, name)
	}

	rv := reflect.TypeOf(d.dc)
	for i := 0; i < rv.NumMethod(); i++ {
		name := changecase.Snake(rv.Method(i).Name)
		fmt.Printf("name: %v\n", name)
		fmt.Printf("hiddenMethods[name]: %v\n", hiddenMethods[name])
		if hiddenMethods[name] {
			continue
		}
		names = append(names, name)
	}

	return names
}

// Freeze implements starlark.Value.
func (d *DrawContext) Freeze() {
	d.frozen = true
}

// Hash implements starlark.Value.
func (d *DrawContext) Hash() (uint32, error) {
	return 0, fmt.Errorf("DrawContext is not hashable")
}

// String implements starlark.Value.
func (d *DrawContext) String() string {
	return fmt.Sprintf("DrawContext(width=%d, height=%d)", d.dc.Width(), d.dc.Height())
}

// Truth implements starlark.Value.
func (d *DrawContext) Truth() starlark.Bool {
	return d != nil
}

// Type implements starlark.Value.
func (d *DrawContext) Type() string {
	return "DrawContext"
}

// format is png or jpeg.
// quality only working for jpeg
func (d *DrawContext) GetBase64(format string, quality int) (_ starlark.String, err error) {
	buf := bytes.NewBuffer(nil)
	switch format {
	case "png":
		err = png.Encode(buf, d.dc.Image())
	case "jpeg":
		err = jpeg.Encode(buf, d.dc.Image(), &jpeg.Options{Quality: quality})
	default:
		return "", fmt.Errorf("not support format: %s", format)
	}
	if err != nil {
		return "", err
	}
	raw := base64.StdEncoding.EncodeToString(buf.Bytes())
	return starlark.String(raw), nil
}

func (d *DrawContext) GetDataURI(format string, quality int) (_ starlark.String, err error) {
	raw, err := d.GetBase64(format, quality)
	if err != nil {
		return "", err
	}
	return starlark.String(fmt.Sprintf("data:image/%s;base64,%s", format, raw.GoString())), nil
}

func NewContext(width, height starlark.Int) (_ *DrawContext, err error) {
	var w, h int
	if err := starlark.AsInt(width, &w); err != nil {
		return nil, err
	}
	if err := starlark.AsInt(height, &h); err != nil {
		return nil, err
	}
	dc := gg.NewContext(w, h)
	dc.SetFontFace(GetFace(16)) // 默认字体
	return &DrawContext{
		dc: dc,
	}, nil
}

var _ starlark.Value = (*DrawContext)(nil)
var _ starlark.HasAttrs = (*DrawContext)(nil)
