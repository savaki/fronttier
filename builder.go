package frontier

import (
	"errors"
	"github.com/savaki/frontier/auth"
	"github.com/savaki/frontier/filter"
)

type BuilderConfig struct {
	routeConfigs []*RouteConfig
	authConfig   *auth.BuilderConfig
	err          error
}

func Builder() *BuilderConfig {
	return &BuilderConfig{}
}

func (self *BuilderConfig) Path(paths ...string) *RouteConfig {
	routeConfig := newRouteBuilder()

	for _, path := range paths {
		routeConfig.Matcher(&PrefixMatcher{path})
	}

	self.routeConfigs = append(self.routeConfigs, routeConfig)
	return routeConfig
}

func (self *BuilderConfig) AuthConfig() *auth.BuilderConfig {
	self.authConfig = auth.Builder()
	return self.authConfig
}

func (self *BuilderConfig) Build() (*Frontier, error) {
	if self.err != nil {
		return nil, self.err
	}

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

	server := &Frontier{routes: routes}
	return server, nil
}
