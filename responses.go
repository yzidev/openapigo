package goas

import "net/http"

// ResponseSpec declares an OpenAPI response for a specific status code.
//
// Schema can be nil to indicate an empty body.
// Description is optional; if empty we'll use http.StatusText(code) where possible.
//
// Example:
//
//	WithResponses(
//	  ResponseSpec{Status: 200, Schema: User{}, Description: "OK"},
//	  ResponseSpec{Status: 401, Schema: ErrorResponse{}, Description: "Unauthorized"},
//	)
type ResponseSpec struct {
	Status      int
	Description string
	Schema      interface{}
}

func (r ResponseSpec) normalizedDescription() string {
	if r.Description != "" {
		return r.Description
	}
	if txt := http.StatusText(r.Status); txt != "" {
		return txt
	}
	return "Response"
}

// WithResponses adds/overrides response entries for a route.
func WithResponses(responses ...ResponseSpec) HandlerOption {
	return func(meta *RouteMeta) {
		meta.Responses = append(meta.Responses, responses...)
	}
}
