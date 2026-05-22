package goas

import "mime/multipart"

// MultipartFile is a marker type for OpenAPI generation.
//
// Use it inside a request schema struct to document multipart/form-data uploads.
// Example:
//
//	type UploadReq struct {
//		File MultipartFile `json:"file"`
//		Note string       `json:"note"`
//	}
//
// When used with WithRequestSchema(UploadReq{}), the builder will render the
// request body as multipart/form-data with a binary file part.
//
// At runtime you still parse files using your framework (r.FormFile, c.FormFile, etc).
// This type is only for spec generation.
type MultipartFile struct{}

// Silence unused import errors when MultipartFile is referenced via multipart.FileHeader.
var _ = multipart.FileHeader{}
