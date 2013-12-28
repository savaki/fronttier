package proxy

import (
	"io"
	"net/http"
	"net/url"
)

func rewrite(target *url.URL, ignoredHeaders []string, req *http.Request) *http.Request {
	outreq := new(http.Request)

	// protocol
	outreq.Proto = "HTTP/1.1"
	outreq.ProtoMajor = 1
	outreq.ProtoMinor = 1
	outreq.Close = false

	path := req.URL.Path
	if len(target.Path) > 1 {
		path = target.Path + path
	}

	// url path
	outreq.URL = &url.URL{
		Scheme:   target.Scheme,
		Opaque:   target.Opaque,
		User:     nil,
		Host:     target.Host,
		Path:     path,
		RawQuery: req.URL.RawQuery,
		Fragment: req.URL.Fragment,
	}

	// http headers
	outreq.Header = make(http.Header)

	for key, values := range req.Header {
		for _, value := range values {
			outreq.Header.Add(key, value)
		}
	}
	outreq.Header.Set("Host", target.Host)

	// content
	outreq.Body = req.Body
	outreq.ContentLength = req.ContentLength

	return outreq
}

func transfer(target http.ResponseWriter, source *http.Response) {
	if source == nil {
		return
	}

	for key, values := range source.Header {
		for _, value := range values {
			target.Header().Add(key, value)
		}
	}

	target.WriteHeader(source.StatusCode)
	io.Copy(target, source.Body)
}
