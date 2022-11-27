// exportable typescript generated from golang
// Copyright (C) 2022  Fabio Prada

package tsrpc

import (
	"regexp"
	"strings"
)

type TSTagJson struct {
	Name      string
	Ignore    bool
	OmitEmpty bool
}

func (t *TSTagJson) parse(tag string) bool {
	*t = TSTagJson{}
	t.Name = ""
	t.Ignore = true
	t.OmitEmpty = false
	re := regexp.MustCompile(`json:\"(.*?)\"`)
	match := re.FindStringSubmatch(tag)
	if len(match) == 0 {
		return false
	}
	t.Name = ""
	s := strings.Split(match[1], ",")
	if len(s) == 1 {
		if s[0] == "omitempty" {
			t.OmitEmpty = true
		}
		t.Ignore = s[0] == "-"
		if s[0] != "" {
			t.Name = s[0]
		}
	}
	if len(s) == 2 {
		t.Ignore = s[0] == "-"
		if s[0] != "" {
			t.Name = s[0]
		}
		if s[1] == "omitempty" {
			t.OmitEmpty = true
		}
	}

	return true
}
