package fronttier

import (
	"net/http"
	"strings"
)

// Matcher provides a standard interface for all sorts of Route
// matchers such as PrefixMatcher and OrMatcher
type Matcher interface {
	Matches(*http.Request) bool
}

// PrefixMatcher is a matcher that matches against the uri prefix
type PrefixMatcher struct {
	prefix string
}

// Match against the URI prefix
func (self *PrefixMatcher) Matches(req *http.Request) bool {
	return strings.HasPrefix(req.URL.Path, self.prefix)
}

// OrMatcher applies a logical Or to all the matchers defined
type OrMatcher struct {
	matchers []Matcher
}

func (self *OrMatcher) Add(matcher Matcher) *OrMatcher {
	self.matchers = append(self.matchers, matcher)
	return self
}

// match if at least one of the matchers matches
func (self *OrMatcher) Matches(req *http.Request) bool {
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

	or := &OrMatcher{}

	for _, matcher := range matchers {
		or.matchers = append(or.matchers, matcher)
	}

	return or
}
