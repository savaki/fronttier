package fronttier

import (
	"github.com/savaki/fronttier/mock"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"testing"
)

func TestRouter(t *testing.T) {
	var router *Router
	var route *Route
	var handler *mock.Handler
	var req *http.Request

	Convey("Given a Router", t, func() {
		router = NewRouter()
		handler = &mock.Handler{}
		req, _ = http.NewRequest("GET", "http://www.acme.com/sample", nil)

		Convey("When I create a #NewRoute", func() {
			route = router.NewRoute()

			Convey("Then I expect the Router to have the new route", func() {
				So(route, ShouldNotBeNil)
				So(len(router.routes), ShouldEqual, 1)
			})
		})

		Convey("When I create a #Route via #HandleFunc", func() {
			router.HandleFunc("/sample", handler.ServeHTTP)

			Convey("Then I expect Router#ServeHTTP to serve to our handler IF path matches", func() {
				router.ServeHTTP(nil, req)

				So(handler.InMethod, ShouldEqual, req.Method)
			})

			Convey("Then I expect Router #ServeHTTP to return 404 if the path doesn't match", func() {
				req, _ = http.NewRequest("GET", "http://www.cnn.com/does-not-match", nil)
				w := &mock.ResponseWriter{}
				router.ServeHTTP(w, req)

				So(handler.InMethod, ShouldEqual, "")
				So(w.StatusCode, ShouldEqual, 404)
			})
		})

		Convey("When I create a #Route via #Handle", func() {
			router.Handle("/sample", handler)

			Convey("Then I expect Router#ServeHTTP to serve to our handler IF path matches", func() {
				router.ServeHTTP(nil, req)

				So(handler.InMethod, ShouldEqual, req.Method)
			})
		})

		Convey("When I create a #Route via #Methods", func() {
			router.Methods("GET").Handler(handler)

			Convey("Then I expect Router#ServeHTTP to serve to our handler IF Method matches", func() {
				router.ServeHTTP(nil, req)

				So(handler.InMethod, ShouldEqual, req.Method)
			})

			Convey("Then I expect Router#ServeHTTP to return 404 to our handler IF Method DOES NOT matches", func() {
				w := &mock.ResponseWriter{}
				req, _ = http.NewRequest("POST", "http://www.acme.com/sample", nil)
				router.ServeHTTP(w, req)

				So(handler.InMethod, ShouldEqual, "")
				So(w.StatusCode, ShouldEqual, 404)
			})
		})

		Convey("When I create a #Route via #PathPrefix", func() {
			router.PathPrefix("/sample").Handler(handler)

			Convey("Then I expect Router#ServeHTTP to serve to our handler IF path matches", func() {
				router.ServeHTTP(nil, req)

				So(handler.InMethod, ShouldEqual, req.Method)
			})
		})
	})
}
