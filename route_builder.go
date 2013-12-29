package fronttier

import (
	"errors"
	"github.com/savaki/fronttier/filter"
	"net/http"
)

type RouteConfig struct {
	sessionFactory bool
	matchers       []Matcher
	filters        []filter.HandlerFilter
	handler        http.Handler
	err            error
}

func newRouteBuilder() *RouteConfig {
	return &RouteConfig{}
}

func (self *RouteConfig) SessionFactory() *RouteConfig {
	self.sessionFactory = true
	return self
}

func (self *RouteConfig) Matcher(matcher Matcher) *RouteConfig {
	self.matchers = append(self.matchers, matcher)
	return self
}

func (self *RouteConfig) Handler(handler http.Handler) *RouteConfig {
	self.handler = handler
	return self
}

func (self *RouteConfig) Filter(filter filter.HandlerFilter) *RouteConfig {
	if self.err == nil {
		self.filters = append(self.filters, filter)
	}
	return self
}

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
