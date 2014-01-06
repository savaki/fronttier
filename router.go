package fronttier

import (
	"net/http"
)

type Router struct {
	routes []*Route
}

func (self *Router) NewRoute() *Route {
	route := &Route{}
	self.routes = append(self.routes, route)
	return route
}

func (self *Router) PathPrefix(prefix string) *Route {
	return self.NewRoute().PathPrefix(prefix)
}

func (self *Router) Methods(methods ...string) *Route {
	return self.NewRoute().Methods(methods...)
}

func (self *Router) HandleFunc(prefix string, handlerFunc http.HandlerFunc) *Router {
	self.NewRoute().PathPrefix(prefix).HandlerFunc(handlerFunc)
	return self
}

func (self *Router) Handle(prefix string, handler http.Handler) *Router {
	return self.HandleFunc(prefix, handler.ServeHTTP)
}

func (self *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for _, route := range self.routes {
		if route.matches(req) {
			route.ServeHTTP(w, req)
			return
		}
	}

	w.WriteHeader(404)
}
