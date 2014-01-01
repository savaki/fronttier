package fronttier

import (
	"errors"
	"github.com/savaki/fronttier/auth"
	"github.com/savaki/fronttier/filter"
)

// BuilderConfig collects all the settings used to construct
// a *Fronttier server.
// To instantiate a new BuilderConfig, call #Builder
// To convert this config into a *Fronttier instance, call #Builder()
type BuilderConfig struct {
	routeConfigs []*RouteConfig
	authConfig   *auth.BuilderConfig
	err          error
}

// The factory method to instantiate new *BuilderConfig.  Start
// here if you want to create a new instance of *Fronttier
func Builder() *BuilderConfig {
	return &BuilderConfig{}
}

// The list of path prefixes to match again.
// Note: / will match everything
func (self *BuilderConfig) Paths(paths ...string) *RouteConfig {
	routeConfig := newRouteBuilder()

	var matchers []Matcher
	for _, path := range paths {
		matchers = append(matchers, &PrefixMatcher{path})
	}
	routeConfig.Matcher(Or(matchers...))

	self.routeConfigs = append(self.routeConfigs, routeConfig)
	return routeConfig
}

// Calling AuthConfig indicates that Fronttier should handle
// authentication.  For more information see the auth package
func (self *BuilderConfig) AuthConfig() *auth.BuilderConfig {
	self.authConfig = auth.Builder()
	return self.authConfig
}

// Instantiate our new Fronttier object
func (self *BuilderConfig) Build() (*Frontier, error) {
	if self.err != nil {
		return nil, self.err
	}

	// build the authentication filters if #AuthConfig was called
	var authFilter filter.HandlerFilter = nil
	var newSessionFilter filter.HandlerFilter = nil
	var err error
	if self.authConfig != nil {
		authFilter, err = self.authConfig.BuildAuthFilter()
		if err != nil {
			return nil, err
		}

		newSessionFilter, _ = self.authConfig.BuildNewSessionFilter()
	}

	// materialize our routes
	var routes []*Route
	for _, routeConfig := range self.routeConfigs {
		if authFilter != nil {
			routeConfig.Filter(authFilter)
		}

		if routeConfig.sessionFactory {
			if newSessionFilter == nil {
				return nil, errors.New("Cannot create Frontier.  #SessionFactory was specified on a Route yet #AuthConfig defined!")
			}
			routeConfig.Filter(newSessionFilter)
		}

		route, err := routeConfig.Build()
		if err != nil {
			return nil, err
		}

		routes = append(routes, route)
	}

	// and lastly, construct our server
	return &Frontier{routes: routes}, nil
}
