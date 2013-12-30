package proxy

import (
	"fmt"
	"github.com/savaki/fronttier/mock"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"testing"
)

func TestService(t *testing.T) {
	var roundTripper *mock.RoundTripper
	var handler http.Handler
	var req *http.Request
	var err error

	Convey("Given a proxy", t, func() {
		roundTripper = &mock.RoundTripper{}
		handler, err = Builder().RoundTripper(roundTripper).Url("http://www.google.com").Build()
		So(err, ShouldBeNil)

		req, _ = http.NewRequest("GET", "http://www.yahoo.com", nil)
		req.RemoteAddr = "10.0.0.2:80"

		Convey("When a request is proxied WITH OUT an X-Forwarded-For header", func() {
			handler.ServeHTTP(nil, req)
			fmt.Printf("%#v\n", roundTripper.Request.Header)
			Convey("Then I expect the X-Forwarded-For header to be set", func() {
				So(roundTripper.Request.Header.Get("X-Forwarded-For"), ShouldNotEqual, "")
			})
		})

		Convey("When a request is proxied WITH an X-Forwarded-For header", func() {
			req.Header.Set("X-Forwarded-For", "10.0.0.1")
			handler.ServeHTTP(nil, req)

			Convey("Then I expect the X-Forwarded-For header to be appended to", func() {
				So(roundTripper.Request.Header.Get("X-Forwarded-For"), ShouldEqual, "10.0.0.1, 10.0.0.2")
			})
		})
	})
}
