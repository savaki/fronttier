package fronttier

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"testing"
)

func mockFilterFunc(headerName, headerValue string) func(http.ResponseWriter, *http.Request, http.HandlerFunc) {
	return func(w http.ResponseWriter, req *http.Request, handler http.HandlerFunc) {
		value := req.Header.Get(headerName)
		if value == "" {
			value = headerValue
		} else {
			value = value + "," + headerValue
		}
		req.Header.Set(headerName, value)
		handler(w, req)
	}
}

func TestFilterFunc(t *testing.T) {
	var filter FilterFunc
	var req *http.Request
	target := func(w http.ResponseWriter, req *http.Request) {
		add(req, "target")
	}

	Convey("Given a Filter", t, func() {
		req, _ = http.NewRequest("GET", "http://www.cnn.com", nil)
		filter = newFilter("before", "after")

		Convey("When I invoke the filter", func() {
			filter(nil, req, target)

			Convey("Then I expect the filter to wrap the target", func() {
				So(req.Header.Get(header), ShouldEqual, "before:target:after")
			})
		})
	})

	Convey("#flatten", t, func() {
		req, _ = http.NewRequest("GET", "http://www.cnn.com", nil)
		var filter1 FilterFunc
		var filter2 FilterFunc
		var all http.HandlerFunc

		Convey("Should apply filters in the order provided", func() {
			filter1 = newFilter("a", "d")
			filter2 = newFilter("b", "c")
			all = flatten([]FilterFunc{filter1, filter2}, target)

			all(nil, req)
			So(req.Header.Get(header), ShouldEqual, "a:b:target:c:d")
		})

		Convey("Should handle when nil filters are provided", func() {
			all = flatten(nil, target)
			all(nil, req)
			So(req.Header.Get(header), ShouldEqual, "target")
		})

		Convey("Should handle when empty filter list is provided", func() {
			all = flatten([]FilterFunc{}, target)
			all(nil, req)
			So(req.Header.Get(header), ShouldEqual, "target")
		})
	})
}
