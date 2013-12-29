package fronttier

import (
	"net/http"
	"strings"
)

type Matcher interface {
	Matches(*http.Request) bool
}

type PrefixMatcher struct {
	prefix string
}

func (self *PrefixMatcher) Matches(req *http.Request) bool {
	return strings.HasPrefix(req.URL.Path, self.prefix)
}
