package fronttier

import (
	"errors"
	"github.com/savaki/fronttier/filter"
	"net/http"
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
	err            error
}

// instantiates a new route builder
func newRouteBuilder() *RouteConfig {
	return &RouteConfig{}
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

// Instantiate a new *Route instance from our configuration.
func (self *RouteConfig) Build() (*Route, error) {
	if self.err != nil {
		return nil, self.err
	} else if self.handler == nil {
		return nil, errors.New("Cannot construct a route without a handler")
	}

	handler := filter.Flatten(self.filters, self.handler)

	return &Route{
		matchers: self.matchers,
		handler:  handler,
	}, nil
}
