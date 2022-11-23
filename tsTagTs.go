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
			t.Type = string(v[len("type="):len(string(v))])
		}
	}
	return true
}
