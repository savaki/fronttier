package mock

import (
	"io/ioutil"
	"net/http"
)

type Handler struct {
	InMethod  string
	InHeader  map[string][]string
	InPath    string
	InContent []byte

	OutStatusCode int
	OutHeader     map[string]string
	OutContent    []byte
}

func (self *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// capture the in
	if req != nil {
		self.InMethod = req.Method
		self.InPath = req.URL.Path
		self.InHeader = make(map[string][]string)

		if req.Body != nil {
			defer req.Body.Close()
			self.InContent, _ = ioutil.ReadAll(req.Body)
		}

		for key, values := range req.Header {
			self.InHeader[key] = values
		}
	}

	// respond with the out
	if w != nil {
		if self.OutHeader != nil {
			for key, value := range self.OutHeader {
				w.Header().Set(key, value)
			}
		}

		statusCode := self.OutStatusCode
		if statusCode == 0 {
			statusCode = 200
		}
		w.WriteHeader(statusCode)

		if self.OutContent != nil {
			w.Write(self.OutContent)
		}
	}
}
