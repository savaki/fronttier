package proxy

import (
	"errors"
	"github.com/savaki/fronttier/mock"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"testing"
)

func TestBuilder(t *testing.T) {
	Convey("#AddFilter", t, func() {
		Convey("add filters to the handler", func() {
			var filterInvoked *bool = new(bool) // check to see filter was invoked
			*filterInvoked = false

			tripper := &mock.RoundTripper{}
			filter := func(req *http.Request, next Handler) (*http.Response, error) {
				*filterInvoked = true
				return nil, errors.New("blah")
			}

			handler, err := Builder().Url("http://www.cnn.com").AddFilter(filter).RoundTripper(tripper).Build()
			So(err, ShouldBeNil)

			// When
			request, _ := http.NewRequest("GET", "http://www.google.com", nil)
			handler.ServeHTTP(nil, request)

			// Then
			So(*filterInvoked, ShouldBeTrue)
			So(tripper.Request, ShouldBeNil)
		})
	})

	Convey("#Build", t, func() {
		Convey("returns error if Url not defined", func() {
			_, err := Builder().Build()

			// Then
			So(err, ShouldNotBeNil)
		})

		Convey("return service if Url is defined", func() {
			service, err := Builder().Url("http://www.google.com").Build()

			So(err, ShouldBeNil)
			So(service, ShouldNotBeNil)
		})

		Convey("sets the Host header for requests that are passed through", func() {
			tripper := mock.RoundTripper{}
			service, _ := Builder().Url("http://www.google.com").RoundTripper(&tripper).Build()
			request, _ := http.NewRequest("GET", "http://www.yahoo.com/sample", nil)

			// When
			service.ServeHTTP(nil, request)

			// Then
			So(tripper.Request.URL.Host, ShouldEqual, "www.google.com")
			So(tripper.Request.Header.Get("Host"), ShouldEqual, "www.google.com")
		})

		Convey("should report errors thrown earlier in the process", func() {
			_, err := Builder().Url(":::this is junk").Build()

			// Then
			So(err, ShouldNotBeNil)
		})

		Convey("the created Handler should return prematurely if the RoundTripper fails", func() {
			tripper := &mock.RoundTripper{
				Request:  nil,
				Response: nil,
				Err:      errors.New("boom!"),
			}
			service, _ := Builder().RoundTripper(tripper).Url("http://www.google.com").Build()
			request, _ := http.NewRequest("GET", "http://www.cnn.com", nil)

			// When
			service.ServeHTTP(nil, request)

			// Then
			// the response writer should not be written to e.g. response writer is nil
		})
	})

	Convey("#NotAuthorizedHandler", t, func() {
		Convey("should invoke the default not authorized handler when required headers not provided", func() {
			notAuthorized := &defaultNotAuthorizedHandler{}
			service, err := Builder().
				Url("http://www.cnn.com").
				RequiredHeaders("X-Foo").
				NotAuthorizedHandler(notAuthorized).
				Build()
			So(err, ShouldBeNil)

			// When
			req, _ := http.NewRequest("GET", "http://www.cnn.com", nil)
			w := &mock.ResponseWriter{}
			service.ServeHTTP(w, req)

			// Then
			So(w.StatusCode, ShouldEqual, 401)
		})
	})

	Convey("#RequiredHeaders", t, func() {
		Convey("should optionally allow authentication to be required by passing a known header", func() {
			service, err := Builder().Url("http://www.cnn.com/").RequiredHeaders("X-Foo").Build()

			// Then
			So(err, ShouldBeNil)
			So(service, ShouldNotBeNil)
		})

		Convey("should assign required headers in the service", func() {
			handler, err := Builder().Url("http://www.cnn.com/").RequiredHeaders("X-Foo").Build()

			// Then
			s := handler.(*proxyService)
			So(err, ShouldBeNil)
			So(len(s.requiredHeaders), ShouldEqual, 1)
			So(s.requiredHeaders, ShouldContain, "X-Foo")
		})

		Convey("should assign a default NotAuthorizedHandler if one wasn't already defined", func() {
			service, err := Builder().
				Url("http://www.cnn.com").
				RequiredHeaders("X-Foo").
				Build()
			So(err, ShouldBeNil)

			// When
			req, _ := http.NewRequest("GET", "http://www.cnn.com", nil)
			w := &mock.ResponseWriter{}
			service.ServeHTTP(w, req)

			// Then
			So(w.StatusCode, ShouldEqual, 401)
		})
	})

	Convey("#Url", t, func() {
		Convey("When I attempt to parse an empty URL", func() {
			_, err := Builder().Url("").Build()

			Convey("Then I expect an error to be thrown", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}
