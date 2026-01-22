package openapi

import (
	"errors"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// ParamType represents a primitive type used for path/query parameters.
type ParamType string

const (
	ParamString  ParamType = "string"
	ParamInteger ParamType = "integer"
	ParamNumber  ParamType = "number"
	ParamBoolean ParamType = "boolean"
)

// QueryParam describes a query parameter for OpenAPI generation.
type QueryParam struct {
	Name        string
	Type        ParamType
	Required    bool
	Description string
}

// WithQueryParams declares query parameters for a route for OpenAPI generation.
func WithQueryParams(params ...QueryParam) HandlerOption {
	return func(meta *RouteMeta) {
		meta.QueryParams = append(meta.QueryParams, params...)
	}
}

type PathParamSpec struct {
	Name        string
	Type        ParamType
	Required    bool
	Description string
}

// WithPathParam declares a typed path parameter (name + primitive type) for OpenAPI generation.
func WithPathParam(name string, typ ParamType, required bool, description string) HandlerOption {
	return func(meta *RouteMeta) {
		meta.PathParams = append(meta.PathParams, PathParamSpec{
			Name:        name,
			Type:        typ,
			Required:    required,
			Description: description,
		})
	}
}

// HeaderParam describes a header parameter for OpenAPI generation.
type HeaderParam struct {
	Name        string
	Type        ParamType
	Required    bool
	Description string
}

// WithHeaderParams declares header parameters for a route for OpenAPI generation.
func WithHeaderParams(params ...HeaderParam) HandlerOption {
	return func(meta *RouteMeta) {
		meta.HeaderParams = append(meta.HeaderParams, params...)
	}
}

func openapiTypeToSchemaType(t ParamType) *openapi3.Types {
	switch t {
	case ParamInteger:
		return &openapi3.Types{"integer"}
	case ParamNumber:
		return &openapi3.Types{"number"}
	case ParamBoolean:
		return &openapi3.Types{"boolean"}
	default:
		return &openapi3.Types{"string"}
	}
}

func addQueryParams(op *openapi3.Operation, qps []QueryParam) {
	for _, qp := range qps {
		if strings.TrimSpace(qp.Name) == "" {
			continue
		}
		p := &openapi3.Parameter{
			Name:        qp.Name,
			In:          openapi3.ParameterInQuery,
			Required:    qp.Required,
			Description: qp.Description,
			Schema:      &openapi3.SchemaRef{Value: &openapi3.Schema{Type: openapiTypeToSchemaType(qp.Type)}},
		}
		op.AddParameter(p)
	}
}

func addHeaderParams(op *openapi3.Operation, hps []HeaderParam) {
	for _, hp := range hps {
		if strings.TrimSpace(hp.Name) == "" {
			continue
		}
		p := &openapi3.Parameter{
			Name:        hp.Name,
			In:          openapi3.ParameterInHeader,
			Required:    hp.Required,
			Description: hp.Description,
			Schema:      &openapi3.SchemaRef{Value: &openapi3.Schema{Type: openapiTypeToSchemaType(hp.Type)}},
		}
		op.AddParameter(p)
	}
}

// parsePrimitive converts a string into a typed primitive (string/int/float/bool).
func parsePrimitive[T any](raw string) (T, error) {
	var z T
	t := reflect.TypeOf(z)
	if t == nil {
		return z, errors.New("nil type")
	}

	switch t.Kind() {
	case reflect.String:
		return any(raw).(T), nil
	case reflect.Bool:
		v, err := strconv.ParseBool(raw)
		if err != nil {
			return z, err
		}
		return any(v).(T), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			return z, err
		}
		out := reflect.New(t).Elem()
		out.SetInt(v)
		return out.Interface().(T), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v, err := strconv.ParseUint(raw, 10, 64)
		if err != nil {
			return z, err
		}
		out := reflect.New(t).Elem()
		out.SetUint(v)
		return out.Interface().(T), nil
	case reflect.Float32, reflect.Float64:
		v, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			return z, err
		}
		out := reflect.New(t).Elem()
		out.SetFloat(v)
		return out.Interface().(T), nil
	default:
		return z, errors.New("unsupported type")
	}
}

// QueryValue reads a query parameter from URL and parses it into the requested type.
func QueryValue[T any](r *http.Request, name string) (T, bool, error) {
	var z T
	if r == nil || r.URL == nil {
		return z, false, errors.New("nil request")
	}
	raw := r.URL.Query().Get(name)
	if raw == "" {
		return z, false, nil
	}
	v, err := parsePrimitive[T](raw)
	if err != nil {
		return z, false, err
	}
	return v, true, nil
}

// QueryValues reads repeated query params (?id=1&id=2) and parses into []T.
func QueryValues[T any](r *http.Request, name string) ([]T, bool, error) {
	if r == nil || r.URL == nil {
		return nil, false, errors.New("nil request")
	}
	raws := r.URL.Query()[name]
	if len(raws) == 0 {
		return nil, false, nil
	}
	out := make([]T, 0, len(raws))
	for _, raw := range raws {
		v, err := parsePrimitive[T](raw)
		if err != nil {
			return nil, true, err
		}
		out = append(out, v)
	}
	return out, true, nil
}

// Utility for tests/examples.
func withQuery(r *http.Request, values url.Values) *http.Request {
	if r.URL == nil {
		return r
	}
	r2 := *r
	u := *r.URL
	u.RawQuery = values.Encode()
	r2.URL = &u
	return &r2
}
