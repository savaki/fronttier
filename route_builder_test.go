package fronttier

import (
	"errors"
	. "github.com/savaki/fronttier/matcher"
	"github.com/savaki/fronttier/mock"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"testing"
)

func TestRouteBuilder(t *testing.T) {
	var builder *RouteConfig
	var err error
	var route *Route
	var handler http.Handler

	Convey("Given a RouteBuilder", t, func() {
		builder = newRouteBuilder()

		Convey("When the builder has an error", func() {
			builder.err = errors.New("I have an error")

			Convey("Then I expect #Build to return an error", func() {
				_, err = builder.Build()

				So(err, ShouldEqual, builder.err)
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When I add a Matcher it appends it to the list of matchers", func() {
			m1 := Prefix("/hello")
			m2 := Prefix("/world")
			builder.Matcher(m1).Matcher(m2)

			Convey("Then I expect the Matchers to have been added in order", func() {
				So(len(builder.matchers), ShouldEqual, 2)
				So(builder.matchers[0], ShouldEqual, m1)
				So(builder.matchers[1], ShouldEqual, m2)
			})

			Convey("The the built Route should have the same matchers", func() {
				handler = &mock.Handler{}
				route, err = builder.Handler(handler).Build()

				So(err, ShouldBeNil)
				So(route, ShouldNotBeNil)
				So(len(route.matchers), ShouldEqual, 2)
				So(route.matchers[0], ShouldResemble, m1)
				So(route.matchers[1], ShouldResemble, m2)
			})
		})

		Convey("When I assign a #Proxy", func() {
			tripper := &mock.RoundTripper{}
			builder.Proxy().RoundTripper(tripper).Url("http://www.eek.com")
			route, err = builder.Matcher(Prefix("/sample")).Build()
			So(err, ShouldBeNil)

			Convey("Then I expect that messages sent to that route to use the proxy", func() {
				w := &mock.ResponseWriter{}
				req, _ := http.NewRequest("GET", "http://www.yahoo.com/sample", nil)
				route.ServeHTTP(w, req)

				So(route.handler, ShouldNotBeNil)
				So(tripper.Request, ShouldNotBeNil)
			})
		})

		Convey("When I assign an incomplete #Proxy", func() {
			builder.Proxy()
			_, err := builder.Matcher(Prefix("/sample")).Build()

			Convey("Then I expect a an err", func() {
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When I assign a Handler", func() {
			handler := &mock.Handler{
				OutHeader: map[string]string{"hello": "world"},
			}
			builder.Handler(handler)

			Convey("Then I expect the built route to have the handler", func() {
				route, err := builder.Handler(handler).Build()

				So(err, ShouldBeNil)
				So(route, ShouldNotBeNil)
				So(route.handler, ShouldResemble, handler)
			})

			Convey("And I also assign a proxy", func() {
				builder.Proxy()

				Convey("Then I expect build to return a HandlerAndProxyErr", func() {
					_, err := builder.Build()

					So(err, ShouldEqual, HandlerAndProxyErr)
				})
			})
		})

		Convey("When I add Filters to a route", func() {
			header := "X-Foo"
			filter1 := &SimpleFilter{header, "1"}
			filter2 := &SimpleFilter{header, "2"}

			builder.Filter(filter1).Filter(filter2)

			Convey("Then I expect the filters to be applied in order", func() {
				handler := &mock.Handler{}
				route, err := builder.Handler(handler).Build()

				So(err, ShouldBeNil)
				So(route, ShouldNotBeNil)

				req, _ := http.NewRequest("GET", "http://www.yahoo.com", nil)
				w := &mock.ResponseWriter{}
				route.ServeHTTP(w, req)
				So(handler.InHeader[header][0], ShouldEqual, "1,2")
			})
		})
	})
}

type SimpleFilter struct {
	header string
	value  string
}

func (self *SimpleFilter) Filter(w http.ResponseWriter, req *http.Request, target http.Handler) {
	value := req.Header.Get(self.header)
	if value == "" {
		req.Header.Set(self.header, self.value)
	} else {
		req.Header.Set(self.header, value+","+self.value)
	}
	target.ServeHTTP(w, req)
}
