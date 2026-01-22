package openapi

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/aizacoders/openapigo/openapi/infer"

	"github.com/getkin/kin-openapi/openapi3"
)

// BuildSpec builds an OpenAPI document from captured routes and config.
func BuildSpec(routes []RouteMeta, cfg Config) *openapi3.T {
	doc := &openapi3.T{
		OpenAPI: "3.0.3",
		Info: &openapi3.Info{
			Title:   cfg.Title,
			Version: cfg.Version,
		},
		Tags:  cfg.Tags,
		Paths: openapi3.NewPaths(),
		Components: &openapi3.Components{
			Schemas:         map[string]*openapi3.SchemaRef{},
			SecuritySchemes: openapi3.SecuritySchemes{},
		},
	}

	// Security schemes
	if cfg.SecuritySchemes != nil {
		for k, v := range cfg.SecuritySchemes {
			doc.Components.SecuritySchemes[k] = v
		}
	}

	// Config-only registered schemas
	if cfg.Schemas != nil {
		for name, sample := range cfg.Schemas {
			if strings.TrimSpace(name) == "" || sample == nil {
				continue
			}
			// Infer schema and then pin it under the provided component name.
			sr := infer.SchemaFrom(doc, sample)
			if doc.Components.Schemas == nil {
				doc.Components.Schemas = map[string]*openapi3.SchemaRef{}
			}
			// If the inferred schema is a $ref, try to copy the referenced schema; otherwise, store Value.
			if sr.Ref != "" {
				// best-effort: if it references an existing component, alias it
				doc.Components.Schemas[name] = &openapi3.SchemaRef{Ref: sr.Ref}
			} else {
				doc.Components.Schemas[name] = sr
			}
		}
	}

	for _, route := range routes {
		path := infer.NormalizePath(route.Path)
		op := &openapi3.Operation{
			Summary:     firstNonEmpty(route.Summary, route.Path),
			Description: route.Description,
			Responses:   &openapi3.Responses{},
		}
		if len(route.Tags) > 0 {
			op.Tags = append(op.Tags, route.Tags...)
		}

		// Path parameters
		if len(route.PathParams) > 0 {
			for _, pp := range route.PathParams {
				if strings.TrimSpace(pp.Name) == "" {
					continue
				}
				op.AddParameter(&openapi3.Parameter{
					Name:        pp.Name,
					In:          openapi3.ParameterInPath,
					Required:    pp.Required,
					Description: pp.Description,
					Schema:      &openapi3.SchemaRef{Value: &openapi3.Schema{Type: openapiTypeToSchemaType(pp.Type)}},
				})
			}
		} else {
			for _, p := range infer.PathParams(path) {
				op.AddParameter(&openapi3.Parameter{
					Name:     p,
					In:       openapi3.ParameterInPath,
					Required: true,
					Schema:   &openapi3.SchemaRef{Value: &openapi3.Schema{Type: &openapi3.Types{"string"}}},
				})
			}
		}

		// Query parameters (declared via WithQueryParams)
		if len(route.QueryParams) > 0 {
			addQueryParams(op, route.QueryParams)
		}

		// Header parameters (declared via WithHeaderParams)
		if len(route.HeaderParams) > 0 {
			addHeaderParams(op, route.HeaderParams)
		}

		if route.RequestSchema != nil {
			schemaRef := infer.RequestSchema(doc, route.RequestSchema)
			op.RequestBody = &openapi3.RequestBodyRef{Value: &openapi3.RequestBody{Required: true, Content: openapi3.NewContentWithJSONSchemaRef(schemaRef)}}
		}

		// Default response behavior (backward-compatible)
		if route.ResponseSchema != nil {
			schemaRef := infer.ResponseSchema(doc, route.ResponseSchema)
			op.Responses.Set("200", &openapi3.ResponseRef{Value: &openapi3.Response{Description: ptr("OK"), Content: openapi3.NewContentWithJSONSchemaRef(schemaRef)}})
		} else {
			op.Responses.Set("200", &openapi3.ResponseRef{Value: &openapi3.Response{Description: ptr("OK")}})
		}

		// Additional/override responses (errors, other success statuses)
		for _, rr := range route.Responses {
			if rr.Status <= 0 {
				continue
			}
			key := strconv.Itoa(rr.Status)
			resp := &openapi3.Response{Description: ptr(rr.normalizedDescription())}
			if rr.Schema != nil {
				schemaRef := infer.ResponseSchema(doc, rr.Schema)
				resp.Content = openapi3.NewContentWithJSONSchemaRef(schemaRef)
			}
			op.Responses.Set(key, &openapi3.ResponseRef{Value: resp})
		}

		// Default error responses (only when user didn't explicitly declare responses)
		// This keeps the "simple usage" experience while still allowing customization.
		if len(route.Responses) == 0 {
			// 401 for secured endpoints
			if route.Security != nil {
				errSchema := infer.ResponseSchema(doc, ErrorResponse{})
				op.Responses.Set("401", &openapi3.ResponseRef{Value: &openapi3.Response{Description: ptr("Unauthorized"), Content: openapi3.NewContentWithJSONSchemaRef(errSchema)}})
			}
			// 400 for likely body-bearing methods
			switch route.Method {
			case http.MethodPost, http.MethodPut, http.MethodPatch:
				errSchema := infer.ResponseSchema(doc, ErrorResponse{})
				op.Responses.Set("400", &openapi3.ResponseRef{Value: &openapi3.Response{Description: ptr("Bad Request"), Content: openapi3.NewContentWithJSONSchemaRef(errSchema)}})
			}
			// default 500
			errSchema := infer.ResponseSchema(doc, ErrorResponse{})
			op.Responses.Set("500", &openapi3.ResponseRef{Value: &openapi3.Response{Description: ptr("Internal Server Error"), Content: openapi3.NewContentWithJSONSchemaRef(errSchema)}})
		}

		if route.Security != nil {
			op.Security = &openapi3.SecurityRequirements{*route.Security}
		}

		// Reuse PathItem if this path already exists, so we don't overwrite other methods.
		item := doc.Paths.Find(path)
		if item == nil {
			item = &openapi3.PathItem{}
		}
		switch route.Method {
		case http.MethodGet:
			item.Get = op
		case http.MethodPost:
			item.Post = op
		case http.MethodPut:
			item.Put = op
		case http.MethodDelete:
			item.Delete = op
		case http.MethodPatch:
			item.Patch = op
		case http.MethodHead:
			item.Head = op
		case http.MethodOptions:
			item.Options = op
		case http.MethodTrace:
			item.Trace = op
		}

		doc.Paths.Set(path, item)
	}

	return doc
}

func firstNonEmpty(v, fallback string) string {
	if v != "" {
		return v
	}
	return fallback
}
