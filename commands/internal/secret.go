package internal

import (
	"fmt"

	"code.cloudfoundry.org/credhub-cli/credhub/credentials"
)

type Secret map[string]interface{}

func NewSecretFromCredential(cred credentials.Credential) (Secret, error) {
	ret := Secret{
		"id":                 cred.Id,
		"name":               cred.Name,
		"version_created_at": cred.VersionCreatedAt,
		"type":               cred.Type,
		"value":              cred.Value,
	}

	return ret, nil
}

func (s Secret) GetSubpath(subpath Subpath) (interface{}, error) {
	return s.getSubpathFor(subpath, map[string]interface{}(s))
}

func (s Secret) getSubpathFor(subpath Subpath, data interface{}) (interface{}, error) {
	outErr := fmt.Errorf("Could not find remainder of subpath `%s'", subpath)
	switch v := data.(type) {
	case int, string, bool:
		if !subpath.Empty() {
			return nil, outErr
		}

		return v, nil

	case map[string]interface{}:
		if subpath.Empty() {
			return v, nil
		}

		allKeys := make([]string, 0, len(v))
		for k := range v {
			allKeys = append(allKeys, k)
		}
		var key string
		key, subpath = subpath.MatchKeyFrom(allKeys...)
		if key == "" {
			return nil, outErr
		}
		return s.getSubpathFor(subpath, v[key])

	case []interface{}:
		if subpath.Empty() {
			return v, nil
		}

		var idx int
		idx, subpath = subpath.MatchIndex(len(v))
		if idx < 0 {
			return nil, outErr
		}

		return s.getSubpathFor(subpath, v[idx])

	default:
		panic("Unknown type!")
	}
}
