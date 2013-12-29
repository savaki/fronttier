package proxy

import (
	"log"
	"net/http"
	"net/url"
)

type proxyService struct {
	target *url.URL
	handle func(req *http.Request) (*http.Response, error)
}

func (self *proxyService) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	outreq := rewrite(self.target, req)

	response, err := self.handle(outreq)
	if err != nil {
		log.Printf("ERROR: remote request failed => %+v\n", err)
		return
	}

	transfer(w, response)
}
