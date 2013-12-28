package frontier

import (
	"net/http"
)

type Frontier struct {
	routes []*Route
}

func (self *Frontier) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// find a router that can handle this request
	for _, route := range self.routes {
		if route.Matches(req) {
			route.ServeHTTP(w, req)
			return
		}
	}
}

func (self *Frontier) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, self)
}
