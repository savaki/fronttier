package matcher

import (
	"net/http"
	"strings"
)

// Matcher provides a standard interface for all sorts of Route
// matchers such as prefixMatcher and orMatcher
type Matcher interface {
	Matches(*http.Request) bool
}

// prefixMatcher is a matcher that matches against the uri prefix
type prefixMatcher struct {
	prefix string
}

// Match against the URI prefix
func (self *prefixMatcher) Matches(req *http.Request) bool {
	return strings.HasPrefix(req.URL.Path, self.prefix)
}

func Prefix(prefix string) Matcher {
	return &prefixMatcher{prefix}
}

// orMatcher applies a logical Or to all the matchers defined
type orMatcher struct {
	matchers []Matcher
}

// match if at least one of the matchers matches
func (self *orMatcher) Matches(req *http.Request) bool {
	for _, matcher := range self.matchers {
		if matcher.Matches(req) {
			return true
		}
	}
	return false
}

// helper method to construct a logical Or match
func Or(matchers ...Matcher) Matcher {
	if len(matchers) == 1 {
		return matchers[0]
	}

	or := &orMatcher{}

	for _, matcher := range matchers {
		or.matchers = append(or.matchers, matcher)
	}

	return or
}
