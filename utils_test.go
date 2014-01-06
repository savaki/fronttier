package fronttier

import (
	"net/http"
)

var header = "X-Header"

func newFilter(before, after string) FilterFunc {
	return func(w http.ResponseWriter, req *http.Request, target http.HandlerFunc) {
		add(req, before)
		target(w, req)
		add(req, after)
	}
}

func add(req *http.Request, value string) {
	current := req.Header.Get(header)
	if current != "" {
		value = current + ":" + value
	}
	req.Header.Set(header, value)
}
