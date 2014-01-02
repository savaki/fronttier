package fronttier

import (
	"errors"
	"github.com/savaki/fronttier/filter"
	. "github.com/savaki/fronttier/matcher"
	"net/http"
)

var (
	NoHandlerDefinedErr  = errors.New("No handler defined")
	HandlerAndBuilderErr = errors.New("Cannot define BOTH a handler AND a builder")
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
	builder        handlerBuilder
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
func (self *RouteConfig) Handler(handler interface{}) *RouteConfig {
	switch v := handler.(type) {
	case http.Handler:
		self.handler = v
	case handlerBuilder:
		self.builder = v
	default:
		self.err = errors.New("#Handler can only accept http.Handler and #handlerBuilder types")
	}
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
	if self.builder != nil {
		return self.builder.Build()

	} else {
		return self.handler, nil
	}
}

// Instantiate a new *Route instance from our configuration.
func (self *RouteConfig) Build() (*Route, error) {
	if self.err != nil {
		return nil, self.err

	} else if self.handler == nil && self.builder == nil {
		return nil, NoHandlerDefinedErr

	} else if self.handler != nil && self.builder != nil {
		return nil, HandlerAndBuilderErr
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
