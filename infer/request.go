package infer

import "github.com/getkin/kin-openapi/openapi3"

func RequestSchema(doc *openapi3.T, sample interface{}) *openapi3.SchemaRef {
	if sr, ok := sample.(*openapi3.SchemaRef); ok && sr != nil {
		return sr
	}
	return SchemaFrom(doc, sample)
}
