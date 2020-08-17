package internal

import "strings"

func CanonizePathForAPI(path string) string {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	return path
}

func CanonizePathForOutput(path string) string {
	return strings.TrimPrefix(path, "/")
}
