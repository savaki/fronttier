package auth

import (
	"github.com/savaki/fronttier/mock"
	"github.com/savaki/fronttier/sessions"
	"net/http"
)

type NewSessionFilter struct {
	idFactory       func() string
	reservedHeaders []string
	template        *http.Cookie
	sessionStore    sessions.Store
}

func (self *NewSessionFilter) createSessionWhenRequired(source http.ResponseWriter, target http.ResponseWriter) {
	reserved := make(map[string]string) // stores our reserved headers

	for _, key := range self.reservedHeaders {
		value := source.Header().Get(key)
		if value != "" {
			reserved[key] = value
		}
	}

	if len(reserved) > 0 {
		session, _ := self.sessionStore.Create(reserved)
		http.SetCookie(target, newCookie(self.template, session.Id))
	}
}

func (self *NewSessionFilter) Filter(w http.ResponseWriter, req *http.Request, target http.Handler) {
	tempWriter := &mock.ResponseWriter{}
	target.ServeHTTP(tempWriter, req)

	copyNonReservedHeaders(self.reservedHeaders, tempWriter, w)
	self.createSessionWhenRequired(tempWriter, w)

	w.WriteHeader(tempWriter.StatusCode)
	tempWriter.WriteTo(w)
}
