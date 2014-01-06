package auth

import (
	"github.com/savaki/fronttier/mock"
	"github.com/savaki/fronttier/sessions"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"testing"
)

func TestAuthFilter(t *testing.T) {
	var filter *AuthFilter
	var request *http.Request
	var handler *mock.Handler
	var w *mock.ResponseWriter
	var sessionStore sessions.Store
	var builder *BuilderConfig

	header := "X-Header"
	values := map[string]string{header: "123"}
	cookieName := "id"
	logoutHeader := "X-Logout"

	Convey("Given an AuthFilter", t, func() {
		handler = &mock.Handler{}
		w = &mock.ResponseWriter{}
		request, _ = http.NewRequest("GET", "http://www.forge.com", nil)
		sessionStore = sessions.Memory()
		builder = Builder().
			ReservedHeaders(header).
			LogoutHeader(logoutHeader).
			CookieName(cookieName).
			SessionStore(sessionStore)

		Convey("When the client attempts to forge a header", func() {
			request.Header.Set(header, "blah")
			filter, _ = builder.BuildAuthFilter()
			filter.Filter(w, request, handler.ServeHTTP)

			Convey("Then the filter should strip out the header", func() {
				So(len(handler.InHeader), ShouldEqual, 0)
			})
		})

		Convey("When the client passes a valid cookie", func() {
			session, _ := sessionStore.Create(values)
			cookie := &http.Cookie{Name: cookieName, Value: session.Id}
			request.AddCookie(cookie)

			filter, _ = builder.SessionStore(sessionStore).BuildAuthFilter()
			filter.Filter(w, request, handler.ServeHTTP)

			Convey("Then values should be retrieved from the sessionStore and placed in the header", func() {
				So(handler.InHeader[header], ShouldContain, "123")
			})
		})

		Convey("When the client passes an INVALID cookie", func() {
			cookie := &http.Cookie{Name: cookieName, Value: "blah"}
			request.AddCookie(cookie)

			filter, _ = builder.SessionStore(sessionStore).BuildAuthFilter()
			filter.Filter(w, request, handler.ServeHTTP)

			Convey("Then the request should be handled normally", func() {
				So(handler.InMethod, ShouldEqual, request.Method)
				So(len(handler.InHeader[header]), ShouldEqual, 0)
			})
		})

		Convey("When the handler returns X-Logout", func() {
			session, _ := sessionStore.Create(values)
			cookie := &http.Cookie{Name: cookieName, Value: session.Id}
			request.AddCookie(cookie)

			handler.OutHeader = map[string]string{logoutHeader: "log-me-out"}

			filter, _ = builder.SessionStore(sessionStore).BuildAuthFilter()
			filter.Filter(w, request, handler.ServeHTTP)

			Convey("Then the session cookie should be cleared", func() {
				session, err := sessionStore.Get(session.Id)

				So(session, ShouldBeNil)
				So(err, ShouldBeNil)
			})
		})

		Convey("When the handler returns []byte content", func() {
			Convey("Then the content should be written to the http.ResponseWriter", func() {
				handler.OutContent = []byte("hello world")

				filter, _ = builder.BuildAuthFilter()
				filter.Filter(w, request, handler.ServeHTTP)

				So(w.String(), ShouldEqual, string(handler.OutContent))
			})
		})
	})
}
