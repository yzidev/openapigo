package infer

import "github.com/getkin/kin-openapi/openapi3"

func ResponseSchema(doc *openapi3.T, sample interface{}) *openapi3.SchemaRef {
	return SchemaFrom(doc, sample)
}
