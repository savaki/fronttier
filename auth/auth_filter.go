package auth

import (
	"github.com/savaki/fronttier/mock"
	"github.com/savaki/fronttier/sessions"
	"net/http"
)

// AuthFilter provides authentication for the services that it's filtering.
type AuthFilter struct {
	idFactory       func() string
	reservedHeaders []string
	logoutHeader    string
	template        *http.Cookie
	sessionStore    sessions.Store
}

func (self *AuthFilter) stripReservedHeaders(req *http.Request) *http.Request {
	for _, key := range self.reservedHeaders {
		req.Header.Del(key)
	}
	return req
}

func (self *AuthFilter) insertSessionInfo(req *http.Request) (*http.Request, *http.Cookie) {
	cookie, _ := req.Cookie(self.template.Name)
	if cookie != nil {
		session, err := self.sessionStore.Get(cookie.Value)
		if err == nil && session != nil {
			for key, value := range session.Values {
				req.Header.Set(key, value)
			}
		}
	}

	return req, cookie
}

func (self *AuthFilter) transferHeaders(cookie *http.Cookie, source http.ResponseWriter, target http.ResponseWriter) {
	copyNonReservedHeaders(self.reservedHeaders, source, target)

	if source.Header()[self.logoutHeader] != nil {
		if cookie != nil {
			self.sessionStore.Delete(cookie.Value)
			http.SetCookie(target, newCookie(self.template, ""))
		}
	}
}

func (self *AuthFilter) Filter(w http.ResponseWriter, req *http.Request, target http.Handler) {
	req = self.stripReservedHeaders(req)
	req, cookie := self.insertSessionInfo(req)

	// capture the response from our service
	tempWriter := &mock.ResponseWriter{}
	target.ServeHTTP(tempWriter, req)

	// and selectively transfer it to the original request
	self.transferHeaders(cookie, tempWriter, w)
	w.WriteHeader(tempWriter.StatusCode)
	tempWriter.WriteTo(w)
}
