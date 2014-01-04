package core

import (
	"errors"
	"github.com/savaki/fronttier/auth"
	"github.com/savaki/fronttier/mock"
	"github.com/savaki/fronttier/sessions"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"testing"
)

func TestBuilder(t *testing.T) {
	var builder *BuilderConfig
	var routeConfig *RouteConfig
	var authConfig *auth.BuilderConfig
	var handler *mock.Handler

	Convey("Given a builder", t, func() {
		builder = Builder()

		Convey("When I call #Build", func() {
			server, err := builder.Build()

			Convey("Then I should get an instance of *Frontier", func() {
				So(err, ShouldBeNil)
				So(server, ShouldNotBeNil)
			})
		})

		Convey("When the builder has an error", func() {
			builder.err = errors.New("I have an error")

			Convey("Then calling #Build should return an error", func() {
				_, err := builder.Build()

				So(err, ShouldNotBeNil)
			})
		})

		Convey("When I call #AuthConfig", func() {
			header := "X-Header"
			authConfig = builder.AuthConfig()
			handler = &mock.Handler{}
			builder.Paths("/a").Handler(handler)

			Convey("Then I expect an error if I have no reserved headers", func() {
				_, err := builder.Build()

				So(err, ShouldNotBeNil)
			})

			Convey("Then I expect all routes to have the AuthFilter", func() {
				authConfig.ReservedHeaders(header)
				server, err := builder.Build()

				So(err, ShouldBeNil)
				So(server, ShouldNotBeNil)

				// Expect reserved headers to be stripped if the AuthFilter was installed
				req, _ := http.NewRequest("GET", "http://www.google.com/a", nil)
				req.Header.Set(header, "blah")
				w := &mock.ResponseWriter{}

				server.ServeHTTP(w, req)

				So(len(handler.InHeader[header]), ShouldEqual, 0)
			})
		})

		Convey("When I call #Paths", func() {
			routeConfig = builder.Paths("/blah/")

			Convey("And I specify the route is a #SessionFactory", func() {
				header := "X-Uid"
				value := "abc123"
				handler = &mock.Handler{OutHeader: map[string]string{header: value}}
				routeConfig.SessionFactory().Handler(handler)

				Convey("Then I expect an error unless the AuthConfig is fulled specified", func() {
					_, err := builder.Build()

					So(err, ShouldNotBeNil)
				})

				Convey("Then I expect a new session to be created if the handler returns the reserved header", func() {
					sessionStore := sessions.Memory()
					builder.AuthConfig().ReservedHeaders(header).SessionStore(sessionStore)
					server, err := builder.Build()

					So(err, ShouldBeNil)
					So(server, ShouldNotBeNil)

					req, _ := http.NewRequest("GET", "http://www.google.com/blah/", nil)
					req.Header.Set("X-Hello", "world")
					w := &mock.ResponseWriter{}

					server.ServeHTTP(w, req)

					So(len(w.Header()["Set-Cookie"]), ShouldEqual, 1)
				})
			})

			Convey("Then I expect an error if I don't define a handler", func() {
				_, err := builder.Build()

				So(err, ShouldNotBeNil)
			})

			Convey("Then I expect requests where DOES match path to be routed to the handler", func() {
				handler := &mock.Handler{}
				routeConfig.Handler(handler)
				server, err := builder.Build()

				req, _ := http.NewRequest("GET", "http://www.yahoo.com/blah/blah", nil)
				w := &mock.ResponseWriter{}
				server.ServeHTTP(w, req)

				So(err, ShouldBeNil)
				So(server, ShouldNotBeNil)
				So(handler.InPath, ShouldEqual, "/blah/blah")
			})

			Convey("Then I expect requests where the path DOES NOT match to NOT be routed to the handler", func() {
				handler := &mock.Handler{}
				routeConfig.Handler(handler)
				server, err := builder.Build()

				req, _ := http.NewRequest("GET", "http://www.yahoo.com/does-not-match", nil)
				w := &mock.ResponseWriter{}
				server.ServeHTTP(w, req)

				So(err, ShouldBeNil)
				So(server, ShouldNotBeNil)
				So(handler.InPath, ShouldEqual, "")
			})
		})
	})
}
