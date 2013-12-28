package proxy

import (
	"bufio"
	"bytes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/savaki/frontier/mock"
	"net/http"
	"net/url"
)

var _ = Describe("Util", func() {
	Context("#rewrite", func() {
		It("should copy header elements", func() {
			target, _ := url.Parse("http://www.google.com/")
			request, _ := http.NewRequest("GET", "http://www.sample.com/", bytes.NewReader([]byte("hello world")))
			request.Header.Set("hello", "world")

			// When
			outreq := rewrite(target, nil, request)

			// Then
			Expect(outreq.Header.Get("hello")).To(Equal("world"), "expected hello == world")
		})

		It("should prepend the path defined in the target to the re-written path", func() {
			target, _ := url.Parse("http://www.google.com/prefix")
			request, _ := http.NewRequest("GET", "http://www.sample.com/abc", nil)

			// When
			outreq := rewrite(target, nil, request)

			// Then
			Expect(outreq.URL.Path).To(Equal("/prefix/abc"), "expected path == /prefix/abc")
		})
	})

	Context("#transfer", func() {
		It("should copy a response to the ResponseWriter", func() {
			request, _ := http.NewRequest("GET", "http://www.sample.com/", bytes.NewReader([]byte("hello world")))
			data := []byte(`HTTP/1.1 200 OK
cache-control: private, max-age=0
content-encoding: gzip
content-type: text/html; charset=UTF-8
content-length: 11

hello world`)
			response, err := http.ReadResponse(bufio.NewReader(bytes.NewReader(data)), request)
			Expect(err).To(BeNil(), "expected no error")

			// When
			w := &mock.ResponseWriter{}
			transfer(w, response)

			// Then
			Expect(w.StatusCode).To(Equal(200), "expected 200 status code")
			Expect(w.Header().Get("content-encoding")).To(Equal("gzip"), "expected content-encoding to be gzip")
			Expect(w.Header().Get("content-length")).To(Equal("11"), "expected content-length to be 11")
			Expect(w.Header().Get("cache-control")).To(Equal("private, max-age=0"), "expected cache-control to be private, max-age=0")
			Expect(string(w.Content)).To(Equal("hello world"), "expected the body to be copied")
		})
	})
})
