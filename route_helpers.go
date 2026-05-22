package goas

import "net/http"

// Req declares a JSON request schema for a route.
func Req(schema any) HandlerOption {
	return WithRequestSchema(schema)
}

// MultipartUpload declares a multipart/form-data request body with a file field.
func MultipartUpload(fileField string, fields ...MultipartField) HandlerOption {
	if fileField == "" {
		fileField = "file"
	}
	return Req(MultipartSchema(fileField, fields...))
}

// Res declares the default success response schema for a route.
func Res(schema any) HandlerOption {
	return WithResponseSchema(schema)
}

// Tags assigns OpenAPI tags to a route.
func Tags(tags ...string) HandlerOption {
	return WithTags(tags...)
}

// Security declares route-level OpenAPI security requirements.
func Security(security *SecurityRequirement) HandlerOption {
	return WithSecurity(security)
}

// Query declares query parameters for a route.
func Query(params ...QueryParam) HandlerOption {
	return WithQueryParams(params...)
}

// Headers declares header parameters for a route.
func Headers(params ...HeaderParam) HandlerOption {
	return WithHeaderParams(params...)
}

// PathParam declares a typed path parameter for a route.
func PathParam(name string, typ ParamType, required bool, description string) HandlerOption {
	return WithPathParam(name, typ, required, description)
}

// Status declares a primary response status for a route.
// When schema is omitted, the route's Res schema is reused if it was set earlier.
func Status(status int, schema ...any) HandlerOption {
	return func(meta *RouteMeta) {
		var s any
		if len(schema) > 0 {
			s = schema[0]
		}
		if s == nil {
			s = meta.ResponseSchema
		}
		meta.Responses = append(meta.Responses, ResponseSpec{Status: status, Schema: s})
	}
}

// Created declares a 201 Created response for a route.
func Created(schema ...any) HandlerOption {
	return Status(http.StatusCreated, schema...)
}

// NoContent declares a 204 No Content response for a route.
func NoContent() HandlerOption {
	return Status(http.StatusNoContent)
}

// Responses declares one or more OpenAPI responses for a route.
func Responses(responses ...ResponseSpec) HandlerOption {
	return WithResponses(responses...)
}

// JSONRoute JSONRouteSpec is a convenience for common JSON APIs.
// It wires request/response schemas + a primary success status code.
//
// You still can override everything by passing explicit options.
//
// Typical usage from adapters:
//
//	r.POST("/users", h, goas.JSONRoute(CreateUser{}, struct{}{}, http.StatusCreated)...)
func JSONRoute(reqSchema any, resSchema any, successStatus int) []HandlerOption {
	opts := make([]HandlerOption, 0, 3)
	if reqSchema != nil {
		opts = append(opts, WithRequestSchema(reqSchema))
	}
	if resSchema != nil {
		opts = append(opts, WithResponseSchema(resSchema))
	}
	if successStatus == 0 {
		successStatus = http.StatusOK
	}
	// declare the primary success response; default errors (400/500/401) are handled by builder when WithResponses isn't used.
	opts = append(opts, WithResponses(ResponseSpec{Status: successStatus, Schema: resSchema}))
	return opts
}
