package page

import "strings"

func ExtractParams(pattern string, path string) map[string]string {
	params := make(map[string]string)
	paths := strings.Split(path, "/")
	for i, v := range strings.Split(pattern, "/") {
		if strings.HasPrefix(v, ":") {
			vt := strings.TrimPrefix(v, ":")
			params[vt] = paths[i]
		}
	}

	return params
}

func BuildPath(pattern string, params map[string]string) string {
	var path string
	patternPaths := strings.Split(pattern, "/")
	for i, pp := range patternPaths {
		if strings.HasPrefix(pp, ":") {
			v, ok := params[strings.TrimPrefix(pp, ":")]
			if !ok {
				continue
			}
			path += v
		} else {
			path += pp
		}

		if i < (len(patternPaths) - 1) {
			path += "/"
		}
	}

	return path
}
