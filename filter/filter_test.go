package filter

import (
	"github.com/savaki/fronttier/mock"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"testing"
)

type SampleFilter struct {
	value string
}

func (self *SampleFilter) Filter(w http.ResponseWriter, req *http.Request, target http.Handler) {
	key := "X-Key"
	value := req.Header.Get(key)
	if len(value) == 0 {
		req.Header.Set(key, self.value)
	} else {
		req.Header.Set(key, value+","+self.value)
	}

	target.ServeHTTP(w, req)
}

func TestFilter(t *testing.T) {
	Convey("#Flatten", t, func() {
		Convey("should process handlers in order", func() {
			f1 := &SampleFilter{"a"}
			f2 := &SampleFilter{"b"}
			target := &mock.Handler{}

			request, _ := http.NewRequest("GET", "http://www.yahoo.com", nil)

			// When
			handler := Flatten([]HandlerFilter{f1, f2}, target)

			// Then
			handler.ServeHTTP(nil, request)
			So(request.Header.Get("X-Key"), ShouldEqual, "a,b")
		})

		Convey("should return the handler if no filters were specified", func() {
			target := &mock.Handler{}

			// When
			handler := Flatten(nil, target)

			// Then
			So(handler, ShouldEqual, target)
		})
	})
}
