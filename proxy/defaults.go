package proxy

import (
	"io"
	"net/http"
)

type defaultNotAuthorizedHandler struct {
}

func (self *defaultNotAuthorizedHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(401)
	io.WriteString(w, "You are not authorized to access this content.")
}
