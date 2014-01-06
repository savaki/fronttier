package fronttier

import (
	"github.com/savaki/fronttier/auth"
	"net/http"
	"sync"
)

type Router struct {
	routes   []*Route
	sessions *auth.BuilderConfig
	mu       *sync.Mutex
	frozen   bool
}

func NewRouter() *Router {
	return &Router{
		mu: &sync.Mutex{},
	}
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

func (self *Router) freeze() {
	if self.mu == nil {
		self.mu = &sync.Mutex{}
	}

	self.mu.Lock()
	defer self.mu.Unlock()

	if self.frozen {
		return
	}

	if self.sessions != nil {
		authFilter, _ := self.sessions.BuildAuthFilter()
		sessionFilter, _ := self.sessions.BuildNewSessionFilter()

		for _, route := range self.routes {
			route.Filter(sessionFilter.Filter)
			if route.sessionFactory {
				route.Filter(authFilter.Filter)
			}
		}
	}
	self.frozen = true
}

func (self *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if !self.frozen {
		self.freeze()
	}

	for _, route := range self.routes {
		if route.matches(req) {
			route.ServeHTTP(w, req)
			return
		}
	}

	w.WriteHeader(404)
}
