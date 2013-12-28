package frontier

import (
	"net/http"
)

type Route struct {
	matchers []Matcher
	handler  http.Handler
}

func (self *Route) Matches(req *http.Request) bool {
	for _, matcher := range self.matchers {
		if !matcher.Matches(req) {
			return false
		}
	}
	return true
}

func (self *Route) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	self.handler.ServeHTTP(w, req)
}
