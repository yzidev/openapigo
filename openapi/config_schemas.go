package openapi

// SchemaRegistry allows registering types that should appear in OpenAPI components/schemas
// without having to attach them to a specific route.
//
// This is the closest Go equivalent of “Spring Boot config-only” schema registration.
// It does NOT automatically infer which route uses which schema; it only ensures
// schemas exist in components for references and tooling.
//
// If you want per-route request/response schemas, keep using JSONRoute(...) or
// WithRequestSchema/WithResponseSchema on that route.
type SchemaRegistry map[string]any
