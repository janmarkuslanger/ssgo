package page

import (
	"errors"
	"strings"
)

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

func BuildPath(pattern string, params map[string]string) (string, error) {
	var path string
	patternPaths := strings.Split(pattern, "/")
	for i, pp := range patternPaths {
		if strings.HasPrefix(pp, ":") {
			v, ok := params[strings.TrimPrefix(pp, ":")]
			if !ok {
				return "", errors.New("could not replace url param: " + pp)
			}
			path += v
		} else {
			path += pp
		}

		if i < (len(patternPaths) - 1) {
			path += "/"
		}
	}

	return path, nil
}
