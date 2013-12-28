package proxy

import (
	"errors"
	"net/http"
	"net/url"
)

// BuilderConfig implements the BuilderConfig pattern to create new proxy instances.
// To obtain a reference to a new BuilderConfig, call proxy.BuilderConfig()
type BuilderConfig struct {
	// what's the target url
	url *url.URL

	// An optional list of headers to require.  If ALL the required headers
	// are not present, then the notAuthorizedHandler gets invoked
	requiredHeaders []string

	// Our underlying network connection.
	roundTripper http.RoundTripper

	// what's the target url
	notAuthorizedHandler http.Handler

	// An optional list of filters to apply prior to invoking our proxy
	// service.
	filters []Filter

	// If one of the BuilderConfig steps causes an error, it will be captured here
	// and returned upon a call to #Build()
	err error
}

// Builder is the starting point for creating new proxy instances.  As
// the name suggests, it uses the builder pattern to construct proxies.
// Once you've finished adding all your options, call #Build to get an
// instance of http.Handler back
func Builder() BuilderConfig {
	return BuilderConfig{}
}

// Url provides a simple method of assigning schema and host in one
// fell shot.  Currently, the path is unused.
func (self BuilderConfig) Url(rawurl string) BuilderConfig {
	if self.err == nil {
		self.url, self.err = url.Parse(rawurl)
	}
	return self
}

func (self BuilderConfig) AddFilter(filter Filter) BuilderConfig {
	if self.err == nil {
		self.filters = append(self.filters, filter)
	}
	return self
}

func (self BuilderConfig) RequiredHeaders(headers ...string) BuilderConfig {
	if self.err == nil {
		self.requiredHeaders = []string(headers)
	}
	return self
}

func (self BuilderConfig) NotAuthorizedHandler(notAuthorizedHandler http.Handler) BuilderConfig {
	if self.err == nil {
		self.notAuthorizedHandler = notAuthorizedHandler
	}
	return self
}

func (self BuilderConfig) RoundTripper(roundTripper http.RoundTripper) BuilderConfig {
	if self.err == nil {
		self.roundTripper = roundTripper
	}
	return self
}

func (self BuilderConfig) Build() (http.Handler, error) {
	if self.err != nil {
		return nil, self.err
	}

	if self.url == nil {
		return nil, errors.New("#Build failed - no target url defined.  Please call #Url before #Build")
	}

	tripper := self.roundTripper
	if tripper == nil {
		tripper = http.DefaultTransport
	}

	notAuthorized := self.notAuthorizedHandler
	if notAuthorized == nil {
		notAuthorized = &defaultNotAuthorizedHandler{}
	}

	handle := flatten(self.filters, tripper.RoundTrip)

	return &proxyService{
		target:          self.url,
		requiredHeaders: self.requiredHeaders,
		handle:          handle,
		notAuthorized:   notAuthorized,
	}, nil
}
