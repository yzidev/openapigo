package infer

import "strings"

// NormalizePath ensures the stored path matches OpenAPI style.
// It converts common router syntaxes:
//
//	/users/:id   -> /users/{id}
//	/users/<id>  -> /users/{id} (optional - not used today)
func NormalizePath(path string) string {
	if path == "" {
		return path
	}
	parts := strings.Split(path, "/")
	for i, p := range parts {
		if len(p) > 1 && strings.HasPrefix(p, ":") {
			name := strings.TrimPrefix(p, ":")
			if name != "" {
				parts[i] = "{" + name + "}"
			}
		}
		// Support <id> style (some routers) as a best-effort.
		if len(p) > 2 && strings.HasPrefix(p, "<") && strings.HasSuffix(p, ">") {
			name := strings.TrimSuffix(strings.TrimPrefix(p, "<"), ">")
			if name != "" {
				parts[i] = "{" + name + "}"
			}
		}
	}
	return strings.Join(parts, "/")
}

// PathParams extracts OpenAPI style path params from /users/{id}.
func PathParams(path string) []string {
	var params []string
	for _, part := range strings.Split(path, "/") {
		if strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}") {
			name := strings.Trim(part, "{}")
			if name != "" {
				params = append(params, name)
			}
		}
	}
	return params
}
