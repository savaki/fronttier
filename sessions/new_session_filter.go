package sessions

import (
	"github.com/savaki/fronttier/mock"
	"net/http"
)

type NewSessionFilter struct {
	idFactory       func() string
	reservedHeaders []string
	template        *http.Cookie
	sessionStore    Store
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

func (self *NewSessionFilter) Filter(w http.ResponseWriter, req *http.Request, handlerFunc http.HandlerFunc) {
	tempWriter := &mock.ResponseWriter{}
	handlerFunc(tempWriter, req)

	copyNonReservedHeaders(self.reservedHeaders, tempWriter, w)
	self.createSessionWhenRequired(tempWriter, w)

	w.WriteHeader(tempWriter.StatusCode)
	tempWriter.WriteTo(w)
}
