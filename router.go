package fronttier

import (
	"github.com/savaki/fronttier/auth"
	"net/http"
)

type Router struct {
	routes   []*Route
	sessions *auth.BuilderConfig
	ready    bool
}

func NewRouter() *Router {
	return &Router{}
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

func (self *Router) Sessions() *auth.BuilderConfig {
	self.sessions = auth.Builder()
	return self.sessions
}

func (self *Router) prepare() {
	if !self.ready {
		if self.sessions != nil {
			authFilter, _ := self.sessions.BuildAuthFilter()
			sessionFilter, _ := self.sessions.BuildNewSessionFilter()

			for _, route := range self.routes {
				route.Filter(authFilter.Filter)
				if route.sessionFactory {
					route.Filter(sessionFilter.Filter)
				}
			}
		}

		self.ready = true
	}
}

func (self *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if !self.ready {
		self.prepare()
	}

	for _, route := range self.routes {
		if route.matches(req) {
			route.ServeHTTP(w, req)
			return
		}
	}

	w.WriteHeader(404)
}
