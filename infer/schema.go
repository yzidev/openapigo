package infer

import (
	"reflect"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// SchemaFrom adds schema components to doc (if needed) and returns a SchemaRef.
//
// sample is typically a zero value of the struct (User{}) or slice ([]User{}).
func SchemaFrom(doc *openapi3.T, sample interface{}) *openapi3.SchemaRef {
	t := reflect.TypeOf(sample)
	if t == nil {
		return &openapi3.SchemaRef{Value: &openapi3.Schema{Type: &openapi3.Types{"object"}}}
	}
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// slices/arrays
	if t.Kind() == reflect.Slice || t.Kind() == reflect.Array {
		itemType := t.Elem()
		if itemType.Kind() == reflect.Ptr {
			itemType = itemType.Elem()
		}
		itemSchema := SchemaFrom(doc, reflect.New(itemType).Elem().Interface())
		return &openapi3.SchemaRef{Value: &openapi3.Schema{Type: &openapi3.Types{"array"}, Items: itemSchema}}
	}

	// primitives
	if t.Kind() != reflect.Struct {
		return &openapi3.SchemaRef{Value: &openapi3.Schema{Type: &openapi3.Types{kindToType(t.Kind())}}}
	}

	name := t.Name()
	if name == "" {
		// anonymous struct
		return &openapi3.SchemaRef{Value: anonymousStructSchema(doc, t)}
	}

	if doc.Components.Schemas == nil {
		doc.Components.Schemas = map[string]*openapi3.SchemaRef{}
	}
	if _, ok := doc.Components.Schemas[name]; !ok {
		doc.Components.Schemas[name] = &openapi3.SchemaRef{Value: namedStructSchema(doc, t)}
	}

	return &openapi3.SchemaRef{Ref: "#/components/schemas/" + name}
}

func anonymousStructSchema(doc *openapi3.T, t reflect.Type) *openapi3.Schema {
	return namedStructSchema(doc, t)
}

func namedStructSchema(doc *openapi3.T, t reflect.Type) *openapi3.Schema {
	s := &openapi3.Schema{Type: &openapi3.Types{"object"}, Properties: map[string]*openapi3.SchemaRef{}}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.Anonymous {
			// ignore for now (could be expanded later)
			continue
		}
		jsonTag := f.Tag.Get("json")
		if jsonTag == "" {
			continue
		}
		parts := strings.Split(jsonTag, ",")
		name := parts[0]
		if name == "" || name == "-" {
			continue
		}

		ft := f.Type
		if ft.Kind() == reflect.Ptr {
			ft = ft.Elem()
		}

		// Special-case: goas.MultipartFile
		if ft.PkgPath() == "github.com/yzidev/goas" && ft.Name() == "MultipartFile" {
			s.Properties[name] = &openapi3.SchemaRef{Value: &openapi3.Schema{Type: &openapi3.Types{"string"}, Format: "binary"}}
			continue
		}

		if ft.Kind() == reflect.Slice || ft.Kind() == reflect.Array {
			it := ft.Elem()
			if it.Kind() == reflect.Ptr {
				it = it.Elem()
			}
			s.Properties[name] = &openapi3.SchemaRef{Value: &openapi3.Schema{Type: &openapi3.Types{"array"}, Items: SchemaFrom(doc, reflect.New(it).Elem().Interface())}}
			continue
		}

		if ft.Kind() == reflect.Struct && ft.PkgPath() != "" {
			s.Properties[name] = SchemaFrom(doc, reflect.New(ft).Elem().Interface())
			continue
		}

		s.Properties[name] = &openapi3.SchemaRef{Value: &openapi3.Schema{Type: &openapi3.Types{kindToType(ft.Kind())}}}
	}
	return s
}

func kindToType(k reflect.Kind) string {
	switch k {
	case reflect.String:
		return "string"
	case reflect.Bool:
		return "boolean"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "integer"
	case reflect.Float32, reflect.Float64:
		return "number"
	default:
		return "object"
	}
}
