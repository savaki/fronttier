package fronttier

import (
	"net/http"
	"strings"
)

type MatcherFunc func(*http.Request) bool

type Route struct {
	matchers []MatcherFunc
	target   http.HandlerFunc
	filters  []FilterFunc
	all      http.HandlerFunc
}

func (self *Route) PathPrefix(prefix string) *Route {
	matcher := func(req *http.Request) bool {
		return strings.HasPrefix(req.URL.Path, prefix)
	}
	return self.Matcher(matcher)
}

func (self *Route) Methods(m ...string) *Route {
	var methods []string
	for _, method := range m {
		methods = append(methods, strings.ToUpper(method))
	}

	matcher := func(req *http.Request) bool {
		for _, method := range methods {
			if method == req.Method {
				return true
			}
		}
		return false
	}

	return self.Matcher(matcher)
}

func (self *Route) Matcher(matcher MatcherFunc) *Route {
	self.matchers = append(self.matchers, matcher)
	return self
}

func (self *Route) Handler(handler http.Handler) *Route {
	return self.HandlerFunc(handler.ServeHTTP)
}

func (self *Route) HandlerFunc(handlerFunc http.HandlerFunc) *Route {
	self.target = handlerFunc
	self.all = flatten(self.filters, self.target)
	return self
}

func (self *Route) Filter(filter FilterFunc) *Route {
	self.filters = append(self.filters, filter)
	self.all = flatten(self.filters, self.target)
	return self
}

func (self *Route) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	self.all(w, req)
}

func (self *Route) matches(req *http.Request) bool {
	for _, matcher := range self.matchers {
		if matcher(req) == false {
			return false
		}
	}

	return true
}
