package fronttier

import (
	"errors"
	. "github.com/savaki/fronttier/matcher"
	"github.com/savaki/fronttier/mock"
	"github.com/savaki/fronttier/proxy"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"testing"
)

func TestRouteBuilder(t *testing.T) {
	var builder *RouteConfig
	var err error
	var r *Route
	var handler http.Handler

	Convey("Given a RouteBuilder", t, func() {
		builder = &RouteConfig{}
		handler = &mock.Handler{}

		Convey("When the builder has an error", func() {
			builder.err = errors.New("I have an error")

			Convey("Then I expect #Build to return an error", func() {
				_, err = builder.Build()

				So(err, ShouldEqual, builder.err)
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When I call #SessionFactory", func() {
			_, err = builder.SessionFactory().Handler(handler).Build()

			Convey("Then I expect #Build to return true for the bool param", func() {
				So(err, ShouldBeNil)
				So(builder.sessionFactory, ShouldBeTrue)
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
				r, err = builder.Handler(handler).Build()

				So(err, ShouldBeNil)
				So(r, ShouldNotBeNil)
				So(len(r.matchers), ShouldEqual, 2)
				So(r.matchers[0], ShouldResemble, m1)
				So(r.matchers[1], ShouldResemble, m2)
			})
		})

		Convey("When I assign a Handler", func() {
			handler := &mock.Handler{
				OutHeader: map[string]string{"hello": "world"},
			}
			builder.Handler(handler)

			Convey("Then I expect the built route to have the handler", func() {
				r, err := builder.Handler(handler).Build()

				So(err, ShouldBeNil)
				So(r, ShouldNotBeNil)
				So(r.handler, ShouldResemble, handler)
			})
		})

		Convey("When I assign an invalid handler via #Handler", func() {
			_, err = builder.Handler("junk").Build()

			Convey("Then err should be HandlerAndBuilderOnlyErr", func() {
				So(err, ShouldEqual, HandlerAndBuilderOnlyErr)
			})
		})

		Convey("When I don't assign either a #Handler or a handlerBuilder", func() {
			_, err = builder.Build()

			Convey("Then err should be NoHandlerDefinedErr", func() {
				So(err, ShouldEqual, NoHandlerDefinedErr)
			})
		})

		Convey("When I assign both a #Handler and a handlerBuilder", func() {
			_, err = builder.Handler(handler).Handler(proxy.Builder()).Build()

			Convey("Then err should be NoHandlerDefinedErr", func() {
				So(err, ShouldEqual, HandlerAndBuilderErr)
			})
		})

		Convey("When I assign a handlerBuilder that throws an error", func() {
			_, err = builder.Handler(proxy.Builder()).Build()

			Convey("Then err should be NoHandlerDefinedErr", func() {
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When I assign a handlerBuilder via #Handler", func() {
			proxyBuilder := proxy.Builder().Url("http://www.eek.com")
			builder.Handler(proxyBuilder)
			handler, err = builder.getHandler()

			Convey("Then I expect the http.Handler to be materialized", func() {
				So(err, ShouldBeNil)
				So(handler, ShouldNotBeNil)
			})
		})

		Convey("When I add Filters to a route", func() {
			header := "X-Foo"
			filter1 := &SimpleFilter{header, "1"}
			filter2 := &SimpleFilter{header, "2"}

			builder.Filter(filter1).Filter(filter2)

			Convey("Then I expect the filters to be applied in order", func() {
				handler := &mock.Handler{}
				r, err := builder.Handler(handler).Build()

				So(err, ShouldBeNil)
				So(r, ShouldNotBeNil)

				req, _ := http.NewRequest("GET", "http://www.yahoo.com", nil)
				w := &mock.ResponseWriter{}
				r.ServeHTTP(w, req)
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
