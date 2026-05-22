package goas

// ErrorResponse is the default error schema used when you don't provide
// explicit response specs via WithResponses.
//
// You can override/replace this by declaring your own responses with
// WithResponses(...) on a route.
type ErrorResponse struct {
	Error string `json:"error"`
}
