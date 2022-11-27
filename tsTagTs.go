// exportable typescript generated from golang
// Copyright (C) 2022  Fabio Prada

package tsrpc

import (
	"regexp"
	"strings"
)

type TSTagTs struct {
	Type   string
	Expand bool
}

func (t *TSTagTs) parse(tag string) bool {
	*t = TSTagTs{Type: "", Expand: false}
	re := regexp.MustCompile(`ts:\"(.*?)\"`)
	match := re.FindStringSubmatch(tag)
	if len(match) == 0 {
		return false
	}
	s := strings.Split(match[1], ",")
	for _, v := range s {
		t.Expand = strings.Trim(v, " ") == "expand"
		if t.Expand {
			return true
		}
		if strings.Contains(v, "type=") {
			t.Type = string(match[1][len("type="):len(string(match[1]))])
		}
	}
	return true
}
