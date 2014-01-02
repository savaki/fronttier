package fronttier

import (
	. "github.com/savaki/fronttier/matcher"
	"net/http"
)

// Route defines a single proxied entity.  For example, I might set
// up /api/foo as a route to my foo service and /api/bar as a route
// to my bar service.
//
// A Route consists of one matchers who ALL must match and a single
// handler.  If you need to apply filters to the handler, you can do
// so using the #Builder
type Route struct {
	matchers []Matcher
	handler  http.Handler
}

// Indicates that this Route can handle the request.  A request can
// only be processed by one route.
func (self *Route) Matches(req *http.Request) bool {
	for _, matcher := range self.matchers {
		if !matcher.Matches(req) {
			return false
		}
	}
	return true
}

// implements the http.Handler method
func (self *Route) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	self.handler.ServeHTTP(w, req)
}
