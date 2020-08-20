package internal

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Subpath string

func NewSubpath(subpath string) Subpath {
	return Subpath(subpath)
}

//returns empty string on no match
func (s Subpath) MatchKeyFrom(keys ...string) (string, Subpath) {
	keys = sort.StringSlice(keys)
	for _, key := range keys {
		if strings.HasPrefix(string(s), key) {
			s = Subpath(strings.TrimPrefix(strings.TrimPrefix(string(s), key), "."))
			return key, s
		}
	}

	return "", s
}

var indexRegex = regexp.MustCompile(`^(?:([0-9]+)|\[([0-9]+)\])(?:\.|$)`)

//returns -1 on no match
func (s Subpath) MatchIndex(length int) (int, Subpath) {
	matches := indexRegex.FindStringSubmatch(string(s))
	if matches == nil {
		return -1, s
	}

	toParse := matches[1]
	if toParse == "" {
		toParse = matches[2]
	}

	ret, err := strconv.Atoi(toParse)
	if err != nil {
		panic(fmt.Sprintf("Could not parse as int: `%s'", toParse))
	}

	if ret >= length {
		return -1, s
	}

	s = Subpath(indexRegex.ReplaceAllString(string(s), ""))
	return ret, s
}

func (s Subpath) Empty() bool {
	return s == ""
}
