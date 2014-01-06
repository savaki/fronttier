package fronttier

import (
	"github.com/savaki/fronttier/mock"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"testing"
)

func TestRoute(t *testing.T) {
	var route *Route
	var req *http.Request
	var handler *mock.Handler

	Convey("Given a Route", t, func() {
		route = &Route{}
		handler = &mock.Handler{}

		Convey("When I call #SessionFactory", func() {
			route.SessionFactory()

			Convey("Then #sessionFactory is set to true", func() {
				So(route.sessionFactory, ShouldBeTrue)
			})
		})

		Convey("When I add a #PathPrefix", func() {
			route.PathPrefix("/sample")

			Convey("Then I expect a matcher to be added to the route", func() {
				So(len(route.matchers), ShouldEqual, 1)
			})

			Convey("Then the route should match by path", func() {
				req, _ = http.NewRequest("GET", "http://www.google.com/sample", nil)
				So(route.matches(req), ShouldBeTrue)

				req, _ = http.NewRequest("GET", "http://www.google.com/no-match", nil)
				So(route.matches(req), ShouldBeFalse)
			})
		})

		Convey("When I add a #Methods matcher", func() {
			route.Methods("GET")

			Convey("Then I expect the matcher to be added to the route", func() {
				So(len(route.matchers), ShouldEqual, 1)
			})

			Convey("Then I expect the matcher to match by Method", func() {
				req, _ = http.NewRequest("GET", "http://www.yahoo.com", nil)
				So(route.matches(req), ShouldBeTrue)
			})

			Convey("Then I expect non-matching methods to return false", func() {
				req, _ = http.NewRequest("PUT", "http://www.yahoo.com", nil)
				So(route.matches(req), ShouldBeFalse)
			})

			Convey("Then I expect the matcher to be case insensitive", func() {
				route = &Route{}
				route.Methods("post").Handler(handler) // need new Route as matchers are logical AND
				req, _ = http.NewRequest("POST", "http://www.yahoo.com", nil)
				So(route.matches(req), ShouldBeTrue)
			})
		})

		Convey("When I add a custom matcher", func() {
			route.Matcher(func(req *http.Request) bool {
				return req.Method == "GET"
			})

			Convey("Then I expect a matcher to be added to the route", func() {
				So(len(route.matchers), ShouldEqual, 1)
			})

			Convey("Then the route should match by path", func() {
				req, _ = http.NewRequest("GET", "http://www.google.com/", nil)
				So(route.matches(req), ShouldBeTrue)

				req, _ = http.NewRequest("NOPE", "http://www.google.com/", nil)
				So(route.matches(req), ShouldBeFalse)
			})
		})

		Convey("When I set Route#target", func() {
			Convey("Via #Handler", func() {
				route.Handler(handler)

				Convey("When I invoke Route#ServeHTTP", func() {
					req, _ = http.NewRequest("GET", "http://www.google.com/", nil)
					route.ServeHTTP(nil, req)

					Convey("Then I expect our handler to be called", func() {
						So(handler.InMethod, ShouldEqual, req.Method)
					})
				})
			})

			Convey("Via #HandlerFunc", func() {
				route.HandlerFunc(handler.ServeHTTP)

				Convey("When I invoke Route#ServeHTTP", func() {
					req, _ = http.NewRequest("GET", "http://www.google.com/", nil)
					route.ServeHTTP(nil, req)

					Convey("Then I expect our handler to be called", func() {
						So(handler.InMethod, ShouldEqual, req.Method)
					})
				})
			})
		})

		Convey("When I add a handler via #Proxy", func() {
			roundTripper := &mock.RoundTripper{}
			route.PathPrefix("/").Proxy("http://www.google.com").ProxyRoundTripper(roundTripper)

			Convey("Then I expect calls to ServeHTTP to route through the proxy", func() {
				route.ServeHTTP(nil, req)

				So(roundTripper.Request.Method, ShouldEqual, req.Method)
				So(roundTripper.Request.Header.Get("Host"), ShouldEqual, "www.google.com")
			})
		})

		Convey("When I add filters via #Filter", func() {
			handlerFunc := func(w http.ResponseWriter, req *http.Request) {
				add(req, "target")
			}
			filter1 := newFilter("a", "d")
			filter2 := newFilter("b", "c")

			route.Filter(filter1)
			route.Filter(filter2)
			route.HandlerFunc(handlerFunc)

			Convey("Then I expect the filters to be applied in order", func() {
				req, _ = http.NewRequest("GET", "http://www.google.com/", nil)
				route.ServeHTTP(nil, req)

				So(req.Header.Get(header), ShouldEqual, "a:b:target:c:d")
			})
		})
	})
}
