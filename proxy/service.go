package proxy

import (
	"log"
	"net/http"
	"net/url"
)

type proxyService struct {
	requiredHeaders []string
	ignoredHeaders  []string
	target          *url.URL
	handle          func(req *http.Request) (*http.Response, error)
	notAuthorized   http.Handler
}

func (self *proxyService) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	outreq := rewrite(self.target, self.ignoredHeaders, req)

	if !self.authorized(req) {
		self.notAuthorized.ServeHTTP(w, req)
		return
	}

	response, err := self.handle(outreq)
	if err != nil {
		log.Printf("ERROR: remote request failed => %+v\n", err)
		return
	}

	transfer(w, response)
}

func (self *proxyService) authorized(req *http.Request) bool {
	if self.requiredHeaders == nil {
		return true
	}

	for _, header := range self.requiredHeaders {
		if req.Header.Get(header) == "" {
			return false
		}
	}

	return true
}
