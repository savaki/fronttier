package sessions

import (
	"github.com/savaki/fronttier/mock"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"testing"
)

func TestAuthFilter(t *testing.T) {
	var filter *AuthFilter
	var req *http.Request
	var handler *mock.Handler
	var w *mock.ResponseWriter
	var sessionStore Store
	var builder *BuilderConfig
	var cookie *http.Cookie
	var valid bool

	header := "X-Header"
	values := map[string]string{header: "123"}
	cookieName := "id"
	logoutHeader := "X-Logout"
	routeName := "sample-route"

	Convey("Given an AuthFilter", t, func() {
		handler = &mock.Handler{}
		w = &mock.ResponseWriter{}
		req, _ = http.NewRequest("GET", "http://www.forge.com", nil)
		sessionStore = Memory()
		builder = Builder().
			ReservedHeaders(header).
			LogoutHeader(logoutHeader).
			CookieName(cookieName).
			SessionStore(sessionStore)

		Convey("When the client attempts to forge a header", func() {
			req.Header.Set(header, "blah")
			filter, _ = builder.BuildAuthFilter()
			filter.Filter(w, req, handler.ServeHTTP, routeName)

			Convey("Then the filter should strip out the header", func() {
				So(len(handler.InHeader), ShouldEqual, 0)
			})
		})

		Convey("When the client passes a valid cookie", func() {
			session, _ := sessionStore.Create(values)
			cookie := &http.Cookie{Name: cookieName, Value: session.Id}
			req.AddCookie(cookie)

			filter, _ = builder.SessionStore(sessionStore).BuildAuthFilter()
			filter.Filter(w, req, handler.ServeHTTP, routeName)

			Convey("Then values should be retrieved from the sessionStore and placed in the header", func() {
				So(handler.InHeader[header], ShouldContain, "123")
			})
		})

		Convey("When the client passes an INVALID cookie", func() {
			cookie := &http.Cookie{Name: cookieName, Value: "blah"}
			req.AddCookie(cookie)

			filter, _ = builder.SessionStore(sessionStore).BuildAuthFilter()
			filter.Filter(w, req, handler.ServeHTTP, routeName)

			Convey("Then the req should be handled normally", func() {
				So(handler.InMethod, ShouldEqual, req.Method)
				So(len(handler.InHeader[header]), ShouldEqual, 0)
			})
		})

		Convey("When the handler returns X-Logout", func() {
			session, _ := sessionStore.Create(values)
			cookie := &http.Cookie{Name: cookieName, Value: session.Id}
			req.AddCookie(cookie)

			handler.OutHeader = map[string]string{logoutHeader: "log-me-out"}

			filter, _ = builder.SessionStore(sessionStore).BuildAuthFilter()
			filter.Filter(w, req, handler.ServeHTTP, routeName)

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
				filter.Filter(w, req, handler.ServeHTTP, routeName)

				So(w.String(), ShouldEqual, string(handler.OutContent))
			})
		})

		Convey("When the #verify func is set AND returns true", func() {
			verifier := &MockSigner{
				OutVerify: true,
			}
			filter, _ := builder.Verifier(verifier).BuildAuthFilter()

			Convey("Then #authorize should always return true", func() {
				cookie, valid = filter.authorize(routeName, req)

				So(cookie, ShouldBeNil)
				So(valid, ShouldBeTrue)
			})
		})

		Convey("When the #sign func is set", func() {
			signer := &MockSigner{}
			filter, _ := builder.Signer(signer).BuildAuthFilter()

			Convey("Then #authorize should call the sign func", func() {
				filter.authorize(routeName, req)

				So(signer.InSign, ShouldResemble, req)
			})
		})
	})
}
