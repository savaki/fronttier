package mock

import (
	"net/http"
)

type ResponseWriter struct {
	header     http.Header
	StatusCode int
	Content    []byte
}

func (self *ResponseWriter) Header() http.Header {
	if self.header == nil {
		self.header = make(http.Header)
	}
	return self.header
}

func (self *ResponseWriter) WriteHeader(statusCode int) {
	self.StatusCode = statusCode
}

func (self *ResponseWriter) Write(data []byte) (int, error) {
	self.Content = append(self.Content, data...)
	return len(data), nil
}
