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
