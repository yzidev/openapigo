package goas

import "github.com/getkin/kin-openapi/openapi3"

// MultipartField defines an extra multipart/form-data field (besides the file).
//
// Type controls the primitive type for that field in the OpenAPI schema.
// If omitted/unknown, it defaults to string.
//
// NOTE: multipart/form-data fields are represented as properties on an object schema.
type MultipartField struct {
	Name        string
	Type        ParamType
	Required    bool
	Description string
}

// MultipartSchema returns a schema sample that produces a multipart/form-data request body.
//
// This returns an *openapi3.SchemaRef directly so we can accurately represent:
// - file as type=string, format=binary
// - primitive fields as string/integer/number/boolean
// - required fields
//
// It is designed for use with the oas builder:
//
//	s.POST("/upload").Req(goas.MultipartSchema("file", goas.MultipartField{Name: "note"})).Res(...)
func MultipartSchema(fileField string, fields ...MultipartField) *openapi3.SchemaRef {
	if fileField == "" {
		fileField = "file"
	}

	s := &openapi3.Schema{Type: &openapi3.Types{"object"}, Properties: map[string]*openapi3.SchemaRef{}, Required: []string{}}

	// file field
	s.Properties[fileField] = &openapi3.SchemaRef{Value: &openapi3.Schema{Type: &openapi3.Types{"string"}, Format: "binary"}}
	s.Required = append(s.Required, fileField)

	// other fields
	for _, f := range fields {
		if f.Name == "" {
			continue
		}
		typ := openapiTypeToSchemaType(f.Type)
		if typ == nil {
			typ = &openapi3.Types{"string"}
		}
		prop := &openapi3.Schema{Type: typ, Description: f.Description}
		s.Properties[f.Name] = &openapi3.SchemaRef{Value: prop}
		if f.Required {
			s.Required = append(s.Required, f.Name)
		}
	}

	return &openapi3.SchemaRef{Value: s}
}
