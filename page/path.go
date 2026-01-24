package page

import (
	"errors"
	"fmt"
	"path"
	"strings"
)

func ExtractParams(pattern string, path string) map[string]string {
	params := make(map[string]string)
	pathParts := splitPathSegments(path)
	for i, v := range splitPathSegments(pattern) {
		if strings.HasPrefix(v, ":") {
			if i >= len(pathParts) {
				continue
			}
			vt := strings.TrimPrefix(v, ":")
			if vt == "" {
				continue
			}
			params[vt] = pathParts[i]
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

func NormalizePagePath(input string) (string, error) {
	cleaned, isRoot, err := normalizePath(input)
	if err != nil {
		return "", err
	}
	if isRoot {
		return "", fmt.Errorf("page path %q resolves to empty output path; use \"index\" or \"index.html\"", input)
	}
	if strings.HasPrefix(cleaned, "..") {
		return "", fmt.Errorf("page path must not traverse outside output dir: %q", input)
	}
	return cleaned, nil
}

func NormalizeRoutePath(input string) (string, error) {
	cleaned, isRoot, err := normalizePath(input)
	if err != nil {
		return "", err
	}
	if strings.HasPrefix(cleaned, "..") {
		return "", fmt.Errorf("page path must not traverse outside output dir: %q", input)
	}
	if isRoot {
		return "/", nil
	}
	return "/" + cleaned, nil
}

func validatePatternMatch(pattern string, path string) error {
	if strings.TrimSpace(pattern) == "" {
		return nil
	}
	patternParts := splitPathSegments(pattern)
	pathParts := splitPathSegments(path)
	if len(patternParts) != len(pathParts) {
		return fmt.Errorf(
			"path %q does not match pattern %q: expected %d segments, got %d",
			path,
			pattern,
			len(patternParts),
			len(pathParts),
		)
	}
	for i, part := range patternParts {
		if strings.HasPrefix(part, ":") {
			if strings.TrimPrefix(part, ":") == "" {
				return fmt.Errorf("pattern %q has empty param name at segment %d", pattern, i+1)
			}
			continue
		}
		if part != pathParts[i] {
			return fmt.Errorf(
				"path %q does not match pattern %q at segment %d: expected %q, got %q",
				path,
				pattern,
				i+1,
				part,
				pathParts[i],
			)
		}
	}
	return nil
}

func splitPathSegments(value string) []string {
	normalized := normalizeForMatch(value)
	if normalized == "" {
		return nil
	}
	return strings.Split(normalized, "/")
}

func normalizeForMatch(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return ""
	}
	trimmed = strings.ReplaceAll(trimmed, "\\", "/")
	trimmed = strings.Trim(trimmed, "/")
	if trimmed == "" {
		return ""
	}
	cleaned := path.Clean(trimmed)
	if cleaned == "." {
		return ""
	}
	return cleaned
}

func normalizePath(input string) (string, bool, error) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return "", false, errors.New("page path must not be empty")
	}
	trimmed = strings.ReplaceAll(trimmed, "\\", "/")
	trimmed = strings.TrimLeft(trimmed, "/")
	if trimmed == "" {
		return "", true, nil
	}
	cleaned := path.Clean(trimmed)
	if cleaned == "." || cleaned == "" {
		return "", true, nil
	}
	return cleaned, false, nil
}
