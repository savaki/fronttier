package fronttier

import (
	"errors"
	"github.com/savaki/fronttier/filter"
	"github.com/savaki/fronttier/proxy"
	"net/http"
)

var (
	NoHandlerDefinedErr = errors.New("No handler defined")
	HandlerAndProxyErr  = errors.New("Cannot define BOTH a handler AND a proxy")
)

// Holds configuration for our route builder.
//
// Note: RouteConfig instances should not be instantiated directly!
// They are created indirectly through #Builder
type RouteConfig struct {
	sessionFactory bool
	matchers       []Matcher
	filters        []filter.HandlerFilter
	handler        http.Handler
	proxyConfig    *proxy.BuilderConfig
	err            error
}

// instantiates a new route builder
func newRouteBuilder() *RouteConfig {
	return &RouteConfig{}
}

func (self *RouteConfig) Proxy() *proxy.BuilderConfig {
	self.proxyConfig = proxy.Builder()
	return self.proxyConfig
}

// Indicates that this Route can create new sessions.  Make to also
// set #AuthConfig in the Fronttier builder to set Fronttier up to
// handle authentication
func (self *RouteConfig) SessionFactory() *RouteConfig {
	self.sessionFactory = true
	return self
}

// Append a new matcher to the list of Matchers for this Route. A
// logical AND will be applied to ALL the matchers to determine if
// this route matches the request
func (self *RouteConfig) Matcher(matcher Matcher) *RouteConfig {
	self.matchers = append(self.matchers, matcher)
	return self
}

// Specify the underlying handler to process the request.  A Route
// has only one handler so calling this again will replace the
// previous Handler
func (self *RouteConfig) Handler(handler http.Handler) *RouteConfig {
	self.handler = handler
	return self
}

// Appends this Filter to the list of Filters to be applied to the
// Handler.  Filters are all AROUND filters and will be executed in the
// order they are applied here.
func (self *RouteConfig) Filter(filter filter.HandlerFilter) *RouteConfig {
	if self.err == nil {
		self.filters = append(self.filters, filter)
	}
	return self
}

func (self *RouteConfig) getHandler() (http.Handler, error) {
	if self.proxyConfig != nil {
		return self.proxyConfig.Build()
	} else {
		return self.handler, nil
	}
}

// Instantiate a new *Route instance from our configuration.
func (self *RouteConfig) Build() (*Route, error) {
	if self.err != nil {
		return nil, self.err

	} else if self.handler == nil && self.proxyConfig == nil {
		return nil, NoHandlerDefinedErr

	} else if self.handler != nil && self.proxyConfig != nil {
		return nil, HandlerAndProxyErr
	}

	handler, err := self.getHandler()
	if err != nil {
		return nil, err
	}

	handler = filter.Flatten(self.filters, handler)

	return &Route{
		matchers: self.matchers,
		handler:  handler,
	}, nil
}
