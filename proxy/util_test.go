package proxy

import (
	"bufio"
	"bytes"
	"github.com/savaki/fronttier/mock"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/url"
	"testing"
)

func TestUtil(t *testing.T) {
	Convey("#rewrite", t, func() {
		Convey("should copy header elements", func() {
			target, _ := url.Parse("http://www.google.com/")
			request, _ := http.NewRequest("GET", "http://www.sample.com/", bytes.NewReader([]byte("hello world")))
			request.Header.Set("hello", "world")

			// When
			outreq := rewrite(target, request)

			// Then
			So(outreq.Header.Get("hello"), ShouldEqual, "world")
		})

		Convey("should prepend the path defined in the target to the re-written path", func() {
			target, _ := url.Parse("http://www.google.com/prefix")
			request, _ := http.NewRequest("GET", "http://www.sample.com/abc", nil)

			// When
			outreq := rewrite(target, request)

			// Then
			So(outreq.URL.Path, ShouldEqual, "/prefix/abc")
		})
	})

	Convey("#transfer", t, func() {
		Convey("should copy a response to the ResponseWriter", func() {
			request, _ := http.NewRequest("GET", "http://www.sample.com/", bytes.NewReader([]byte("hello world")))
			data := []byte(`HTTP/1.1 200 OK
cache-control: private, max-age=0
content-encoding: gzip
content-type: text/html; charset=UTF-8
content-length: 11

hello world`)
			response, err := http.ReadResponse(bufio.NewReader(bytes.NewReader(data)), request)
			So(err, ShouldBeNil)

			// When
			w := &mock.ResponseWriter{}
			transfer(w, response)

			// Then
			So(w.StatusCode, ShouldEqual, 200)
			So(w.Header().Get("content-encoding"), ShouldEqual, "gzip")
			So(w.Header().Get("content-length"), ShouldEqual, "11")
			So(w.Header().Get("cache-control"), ShouldEqual, "private, max-age=0")
			So(string(w.Content), ShouldEqual, "hello world")
		})
	})
}
