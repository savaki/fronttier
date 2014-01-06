package proxy

import (
	"github.com/savaki/fronttier/mock"
	"net"
	"net/http"
	"net/url"
)

// Hop-by-hop headers. These are removed when sent to the backend.
// http://www.w3.org/Protocols/rfc2616/rfc2616-sec13.html
var hopHeaders = []string{
	"Connection",
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Te", // canonicalized version of "TE"
	"Trailers",
	"Transfer-Encoding",
	"Upgrade",
}

func rewrite(target *url.URL, req *http.Request) *http.Request {
	outreq := new(http.Request)

	outreq.Method = req.Method

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

Loop:
	for key, values := range req.Header {
		for _, hopHeader := range hopHeaders {
			if key == hopHeader {
				continue Loop
			}
		}

		for _, value := range values {
			outreq.Header.Add(key, value)
		}
	}
	outreq.Header.Set("Host", target.Host)

	// x-forwarded-for
	if host, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
		if prior := outreq.Header.Get("X-Forwarded-For"); prior != "" {
			host = prior + ", " + host
		}
		outreq.Header.Set("X-Forwarded-For", host)
	}

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
	mock.WriteTo(target, source.Body)
}
