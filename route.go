package fronttier

import (
	"github.com/savaki/fronttier/proxy"
	"net/http"
	"strings"
)

type MatcherFunc func(*http.Request) bool

type Route struct {
	matchers       []MatcherFunc
	target         http.HandlerFunc
	filters        []FilterFunc
	proxyConfig    *proxy.BuilderConfig
	all            http.HandlerFunc
	sessionFactory bool
	err            error
	name           string
}

func (self *Route) PathPrefix(prefix string) *Route {
	matcher := func(req *http.Request) bool {
		return strings.HasPrefix(req.URL.Path, prefix)
	}
	return self.Matcher(matcher)
}

func (self *Route) Name(name string) *Route {
	self.name = name
	return self
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

func (self *Route) SessionFactory() *Route {
	self.sessionFactory = true
	return self
}

func (self *Route) withProxy(f func()) *Route {
	if self.proxyConfig == nil {
		self.proxyConfig = proxy.Builder()
	}
	f()
	handler, err := self.proxyConfig.Build()
	if err != nil {
		self.err = err
	} else {
		self.target = handler.ServeHTTP
		self.all = flatten(self.filters, self.target)
	}
	return self
}

func (self *Route) Err() error {
	return self.err
}

func (self *Route) Proxy(rawurl string) *Route {
	return self.withProxy(func() {
		self.proxyConfig.Url(rawurl)
	})
}

func (self *Route) ProxyRoundTripper(roundTripper http.RoundTripper) *Route {
	return self.withProxy(func() {
		self.proxyConfig.RoundTripper(roundTripper)
	})
}

func (self *Route) Filter(filter FilterFunc) *Route {
	self.filters = append(self.filters, filter)
	self.all = flatten(self.filters, self.target)
	return self
}

// works just like #Filter except the filter will be provided with the name of
// this route.  If no route name has been specified, name will be the empty string
func (self *Route) FilterWithName(filterWithName FilterWithNameFunc) *Route {
	filter := func(w http.ResponseWriter, req *http.Request, handler http.HandlerFunc) {
		filterWithName(w, req, handler, self.name)
	}
	return self.Filter(filter)
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
