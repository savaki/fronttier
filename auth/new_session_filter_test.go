package auth

import (
	"errors"
	"github.com/savaki/fronttier/mock"
	"github.com/savaki/fronttier/sessions"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"testing"
)

func TestNewSessionFilter(t *testing.T) {
	var builder *BuilderConfig
	var handler *mock.Handler
	var sessionStore sessions.Store
	header := "X-Header"

	Convey("Given a Builder", t, func() {
		builder = Builder().ReservedHeaders(header)

		Convey("When the builder has an error", func() {
			builder.err = errors.New("I have an error")

			Convey("Then #Build should return that error", func() {
				_, err := builder.BuildNewSessionFilter()

				So(err, ShouldNotBeNil)
				So(err, ShouldEqual, builder.err)
			})
		})

		Convey("When I call build", func() {
			filter, err := builder.BuildNewSessionFilter()

			Convey("Then all the properties are assigned", func() {
				So(err, ShouldBeNil)
				So(filter, ShouldNotBeNil)

				So(filter.idFactory, ShouldNotBeNil)
				So(filter.template, ShouldNotBeNil)
				So(filter.sessionStore, ShouldNotBeNil)
				So(filter.reservedHeaders, ShouldNotBeNil)
				So(len(filter.reservedHeaders), ShouldBeGreaterThan, 0)
			})
		})
	})

	Convey("Given a NewSessionFilter", t, func() {
		value := "123"
		handler = &mock.Handler{OutHeader: map[string]string{header: value}}
		sessionStore = sessions.Memory()

		Convey("When the handler returns the reserved headers, then I expect to create a new session", func() {
			filter, err := Builder().
				SessionStore(sessionStore).
				ReservedHeaders(header).
				BuildNewSessionFilter()

			So(err, ShouldBeNil)
			So(filter, ShouldNotBeNil)

			req, _ := http.NewRequest("GET", "http://www.google.com", nil)
			w := &mock.ResponseWriter{}

			filter.Filter(w, req, handler)

			So(w.Header()["Set-Cookie"], ShouldNotBeNil)
			So(len(w.Header()["Set-Cookie"]), ShouldEqual, 1)
		})

		Convey("When the client sends content", func() {
			Convey("Then the content should be passed through the filter to the handler", func() {
				filter, err := Builder().ReservedHeaders(header).BuildNewSessionFilter()

				So(err, ShouldBeNil)
				So(filter, ShouldNotBeNil)

				content := "hello world"
				handler.OutContent = []byte(content)
				req, _ := http.NewRequest("GET", "http://www.google.com", nil)
				w := &mock.ResponseWriter{}

				filter.Filter(w, req, handler)

				So(string(w.Bytes()), ShouldEqual, content)
			})
		})
	})
}
