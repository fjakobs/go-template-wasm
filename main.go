package main

import (
	"bytes"
	"text/template"
	"syscall/js"
)

// This calls a JS function from Go.
func main() {
	c := make(chan int)

	js.Global().Set("template", js.FuncOf(TemplateJs))

	<-c
}

// Convert js.Value to a Go native type (map, slice, number, bool, or string)
func jsValueToGoType(value js.Value) any {
	switch value.Type() {
	case js.TypeUndefined, js.TypeNull:
		return nil
	case js.TypeBoolean:
		return value.Bool()
	case js.TypeNumber:
		return value.Float()
	case js.TypeString:
		return value.String()
	case js.TypeFunction:
		return func(args ...interface{}) (interface{}, error) {
			return value.Invoke(args...), nil
		}
	case js.TypeObject:
		// Handle Arrays separately
		if js.Global().Get("Array").Call("isArray", value).Bool() {
			length := value.Length()
			slice := make([]any, length)
			for i := 0; i < length; i++ {
				slice[i] = jsValueToGoType(value.Index(i))
			}
			return slice
		}
		// Handle generic objects
		return jsValueToMap(value)
	default:
		return value // Return raw js.Value as fallback
	}
}

// Convert js.Value to map[string]any recursively
func jsValueToMap(v js.Value) map[string]any {
	result := make(map[string]any)
	keys := js.Global().Get("Object").Call("keys", v) // Get object keys

	// Iterate over keys and fetch values
	length := keys.Length()
	for i := 0; i < length; i++ {
		key := keys.Index(i).String()
		result[key] = jsValueToGoType(v.Get(key)) // Recursively process value
	}
	return result
}

func TemplateJs(this js.Value, inputs []js.Value) interface{} {
	values := make(map[string]any)
	values = jsValueToMap(inputs[0])
	templateStr := inputs[1].String()
	jsFuncs := jsValueToMap(inputs[2])

	str, _ := Template(values, templateStr, jsFuncs)
	return str
}

func Template(values map[string]any, templateStr string, funcs template.FuncMap) (string, error) {
	writer := &bytes.Buffer{}
	tmpl, err := template.New("").Funcs(funcs).Parse(templateStr)
	if err != nil {
		return "", err
	}

	tmpl.Execute(writer, values)
	return writer.String(), nil
}