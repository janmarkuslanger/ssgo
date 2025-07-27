package page

import "strings"

func ExtractPattern(pattern string, path string) map[string]string {
	params := make(map[string]string)
	paths := strings.Split(path, "/")
	for i, v := range strings.Split(pattern, "/") {
		if !strings.HasPrefix(v, ":") {
			continue
		}

		vt := strings.TrimPrefix(v, ":")
		params[vt] = paths[i]

	}

	return params
}
