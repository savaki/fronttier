package fronttier

import (
	"github.com/savaki/fronttier/sessions"
	"net/http"
)

type Router struct {
	filters  []FilterFunc
	routes   []*Route
	sessions *sessions.BuilderConfig
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

func (self *Router) Filter(filter FilterFunc) *Router {
	self.filters = append(self.filters, filter)
	return self
}

func (self *Router) HandleFunc(prefix string, handlerFunc http.HandlerFunc) *Router {
	self.NewRoute().PathPrefix(prefix).HandlerFunc(handlerFunc)
	return self
}

func (self *Router) Handle(prefix string, handler http.Handler) *Router {
	return self.HandleFunc(prefix, handler.ServeHTTP)
}

func (self *Router) Sessions() *sessions.BuilderConfig {
	self.sessions = sessions.Builder()
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

		for _, route := range self.routes {
			for _, filter := range self.filters {
				route.Filter(filter)
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
