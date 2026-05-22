package goas

import "github.com/getkin/kin-openapi/openapi3"

// routeIsMultipartUpload attempts to detect multipart/form-data request bodies.
// We treat any object schema that contains a property with format=binary as a multipart upload.
func routeIsMultipartUpload(schemaRef *openapi3.SchemaRef) bool {
	if schemaRef == nil {
		return false
	}

	// We expect callers to resolve refs (if any) before calling this.
	if schemaRef.Ref != "" {
		return false
	}

	s := schemaRef.Value
	if s == nil {
		return false
	}
	// If an object has any binary property, we consider it multipart.
	if s.Properties != nil {
		for _, p := range s.Properties {
			if p == nil {
				continue
			}
			if p.Value != nil && p.Value.Format == "binary" {
				return true
			}
			// shallow nested
			if p.Value != nil && p.Value.Properties != nil {
				for _, np := range p.Value.Properties {
					if np != nil && np.Value != nil && np.Value.Format == "binary" {
						return true
					}
				}
			}
		}
	}
	return false
}
