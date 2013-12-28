package mock

import (
	"net/http"
)

type RoundTripper struct {
	Request  *http.Request
	Response *http.Response
	Err      error
}

func (self *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	self.Request = req
	return self.Response, self.Err
}
