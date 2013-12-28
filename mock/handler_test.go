package mock

import (
	"bytes"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"testing"
)

func TestHandler(t *testing.T) {
	Convey("#ServeHTTP", t, func() {
		Convey("captures the request path", func() {
			request, _ := http.NewRequest("GET", "http://www.google.com/mail/", nil)

			// When
			handler := &Handler{}
			handler.ServeHTTP(nil, request)

			// Then
			So(handler.InPath, ShouldEqual, "/mail/")
		})

		Convey("captures the request content", func() {
			content := "hello world"
			request, _ := http.NewRequest("GET", "http://www.google.com/mail/", bytes.NewReader([]byte(content)))

			// When
			handler := &Handler{}
			handler.ServeHTTP(nil, request)

			// Then
			So(string(handler.InContent), ShouldEqual, content)
		})

		Convey("expected headers and values to be saved by handler", func() {
			request, _ := http.NewRequest("GET", "http://www.google.com/mail/", nil)
			request.Header.Set("hello", "world")

			// When
			handler := &Handler{}
			handler.ServeHTTP(nil, request)

			// Then
			So(handler.InMethod, ShouldEqual, "GET")
			So(len(handler.InHeader["Hello"]), ShouldEqual, 1)
			So(handler.InHeader["Hello"], ShouldContain, "world")
		})

		Convey("write the specified status code", func() {
			handler := &Handler{OutStatusCode: 201}
			w := &ResponseWriter{}

			// When
			handler.ServeHTTP(w, nil)

			// Then
			So(w.StatusCode, ShouldEqual, handler.OutStatusCode)
		})

		Convey("defaults status code to 200", func() {
			handler := &Handler{}
			w := &ResponseWriter{}

			// When
			handler.ServeHTTP(w, nil)

			// Then
			So(w.StatusCode, ShouldEqual, 200)
		})

		Convey("returns the out headers", func() {
			header := "X-Hello"
			handler := &Handler{OutHeader: map[string]string{header: "world"}}
			w := &ResponseWriter{}

			// When
			handler.ServeHTTP(w, nil)

			// Then
			So(w.Header().Get(header), ShouldEqual, handler.OutHeader[header])
		})

		Convey("returns the out content", func() {
			content := "hello world"
			handler := &Handler{OutContent: []byte(content)}
			w := &ResponseWriter{}

			// When
			handler.ServeHTTP(w, nil)

			// Then
			So(string(w.Content), ShouldEqual, content)
		})
	})
}
