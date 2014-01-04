package core

import (
	"net/http"
)

// Fronttier is an authenticating reverse-proxy design to act as the
// first tier of a service oriented architecture.  The goal is to provide a
// simple mechanism for (a) services be added to a site and (b) authentication
type Frontier struct {
	routes []*Route
}

// Fronttier is usable as a standard http.Handler
func (self *Frontier) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// find a router that can handle this request
	for _, route := range self.routes {
		if route.Matches(req) {
			route.ServeHTTP(w, req)
			return
		}
	}
}

// Alternately, you can start Fronttier directly
func (self *Frontier) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, self)
}
