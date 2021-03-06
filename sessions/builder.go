package sessions

import (
	"errors"
	"net/http"
)

type BuilderConfig struct {
	cookieName        string
	cookiePath        string
	cookieDomain      string
	idFactory         func() string
	timeoutMinutes    int
	timeoutMinutesSet bool
	logoutHeader      string
	reservedHeaders   []string
	sessionStore      Store
	sign              SignerFunc
	verify            VerifierFunc
	err               error
}

func Builder() *BuilderConfig {
	return &BuilderConfig{}
}

func (self *BuilderConfig) CookieName(cookieName string) *BuilderConfig {
	self.cookieName = cookieName
	return self
}

func (self *BuilderConfig) CookiePath(cookiePath string) *BuilderConfig {
	self.cookiePath = cookiePath
	return self
}

func (self *BuilderConfig) CookieDomain(cookieDomain string) *BuilderConfig {
	self.cookieDomain = cookieDomain
	return self
}

func (self *BuilderConfig) IdFactory(idFactory func() string) *BuilderConfig {
	self.idFactory = idFactory
	return self
}

func (self *BuilderConfig) TimeoutMinutes(timeoutMinutes int) *BuilderConfig {
	self.timeoutMinutes = timeoutMinutes
	self.timeoutMinutesSet = true
	return self
}

func (self *BuilderConfig) ReservedHeaders(reservedHeaders ...string) *BuilderConfig {
	self.reservedHeaders = reservedHeaders
	return self
}

func (self *BuilderConfig) LogoutHeader(logoutHeader string) *BuilderConfig {
	self.logoutHeader = logoutHeader
	return self
}

func (self *BuilderConfig) SessionStore(sessionStore Store) *BuilderConfig {
	self.sessionStore = sessionStore
	return self
}

func (self *BuilderConfig) Signer(signer Signer) *BuilderConfig {
	return self.SignerFunc(signer.Sign)
}

func (self *BuilderConfig) SignerFunc(sign SignerFunc) *BuilderConfig {
	self.sign = sign
	return self
}

func (self *BuilderConfig) Verifier(verifier Verifier) *BuilderConfig {
	return self.VerifierFunc(verifier.Verify)
}

func (self *BuilderConfig) VerifierFunc(verify VerifierFunc) *BuilderConfig {
	self.verify = verify
	return self
}

func (self *BuilderConfig) toCookieTemplate() *http.Cookie {
	cookieName := self.cookieName
	if cookieName == "" {
		cookieName = "id"
	}

	path := self.cookiePath
	if path == "" {
		path = "/"
	}

	maxAge := self.timeoutMinutes * 60
	if !self.timeoutMinutesSet {
		maxAge = 15 * 60
	}

	return &http.Cookie{
		Name:   cookieName,
		Path:   path,
		Domain: self.cookieDomain,
		MaxAge: maxAge,
	}
}

func (self *BuilderConfig) toIdFactory() func() string {
	idFactory := self.idFactory
	if idFactory == nil {
		idFactory = defaultIdFactory
	}

	return idFactory
}

func (self *BuilderConfig) toSessionStore() Store {
	sessionStore := self.sessionStore
	if sessionStore == nil {
		sessionStore = defaultSessionStore
	}

	return sessionStore
}

func (self *BuilderConfig) toLogoutHeader() string {
	logoutHeader := self.logoutHeader
	if logoutHeader == "" {
		logoutHeader = "X-Logout"
	}

	return logoutHeader
}

func (self *BuilderConfig) validate() error {
	if self.err != nil {
		return self.err

	} else if self.reservedHeaders == nil || len(self.reservedHeaders) == 0 {
		return errors.New("AuthFilter requires at least one reserved header!")
	}

	return nil
}

func (self *BuilderConfig) BuildAuthFilter() (*AuthFilter, error) {
	if err := self.validate(); err != nil {
		return nil, err
	}

	return &AuthFilter{
		reservedHeaders: self.reservedHeaders,
		logoutHeader:    self.toLogoutHeader(),
		sessionStore:    self.toSessionStore(),
		idFactory:       self.toIdFactory(),
		template:        self.toCookieTemplate(),
		sign:            self.sign,
		verify:          self.verify,
	}, nil
}

func (self *BuilderConfig) BuildNewSessionFilter() (*NewSessionFilter, error) {
	if err := self.validate(); err != nil {
		return nil, err
	}

	return &NewSessionFilter{
		reservedHeaders: self.reservedHeaders,
		sessionStore:    self.toSessionStore(),
		idFactory:       self.toIdFactory(),
		template:        self.toCookieTemplate(),
	}, nil
}
