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

func SplitPath(path string) (string, string) {
	parts := strings.Split(path, ":")
	if len(parts) == 1 {
		return parts[0], ""
	}

	return strings.Join(parts[:len(parts)-1], ":"), parts[len(parts)-1]
}

const (
	CredTypeValue       = "value"
	CredTypeUser        = "user"
	CredTypePassword    = "password"
	CredTypeCertificate = "certificate"
	CredTypeRSA         = "rsa"
	CredTypeSSH         = "ssh"
	CredTypeJSON        = "json"
)
