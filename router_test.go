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

		Convey("When indicate session management via #Sessions", func() {
			reserved := "X-Reserved"
			router.Sessions().ReservedHeaders(reserved)
			router.NewRoute().Handler(handler).SessionFactory()

			Convey("Then I expect to server requests through the handler", func() {
				// set the OutHeader so that a new session gets created
				handler.OutHeader = map[string]string{reserved: "blah"}
				handler.OutStatusCode = 200

				// attempt to forge input header so we can verify its stripped
				req.Header.Add(reserved, "ping")

				w := &mock.ResponseWriter{}
				router.ServeHTTP(w, req)

				// verify session created
				So(w.Header().Get("Set-Cookie"), ShouldNotEqual, "")

				// verify header stripped
				So(len(handler.InHeader[reserved]), ShouldEqual, 0)
			})
		})

		Convey("When I add a global #Filter", func() {
			filter := func(w http.ResponseWriter, req *http.Request, handlerFunc http.HandlerFunc) {

			}
			router.Filter(filter)
		})

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

		Convey("When I apply global filters via #Filter", func() {
			headerName := "X-Sample"
			filter1 := mockFilterFunc(headerName, "a")
			filter2 := mockFilterFunc(headerName, "b")
			router.Filter(filter1).Filter(filter2)
			router.Methods("GET").Handler(handler)

			Convey("Then I expect those filters to be executed in the order provided", func() {
				router.ServeHTTP(nil, req)

				So(handler.InReq, ShouldNotBeNil)
				So(handler.InReq.Header.Get(headerName), ShouldEqual, "a,b")
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
