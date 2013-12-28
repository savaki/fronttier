package proxy

import (
	"errors"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"testing"
)

func TestFilter(t *testing.T) {
	Convey("#flatten", t, func() {
		Convey("should return the target if no filtering needs to happen", func() {
			gripe := errors.New("blah")
			var filters []Filter
			var target = func(*http.Request) (*http.Response, error) {
				return nil, gripe
			}

			// When
			result := flatten(filters, target)

			// Then
			response, err := result(nil)
			So(response, ShouldBeNil)
			So(err, ShouldEqual, gripe)
		})

		Convey("should process filters before the target", func() {
			gripe := errors.New("blah")
			filter := func(req *http.Request, next Handler) (*http.Response, error) {
				return nil, gripe
			}
			filters := []Filter{filter}
			targetInvoked := false
			target := func(*http.Request) (*http.Response, error) {
				targetInvoked = true
				return nil, nil
			}

			// When
			handler := flatten(filters, target)
			handler(nil)

			So(targetInvoked, ShouldBeFalse)
		})

		Convey("should process the FIRST filter provided FIRST", func() {
			gripe := errors.New("blah")
			filter1 := func(req *http.Request, next Handler) (*http.Response, error) {
				return nil, gripe
			}

			filter2Invoked := false
			filter2 := func(req *http.Request, next Handler) (*http.Response, error) {
				filter2Invoked = true
				return nil, nil
			}

			targetInvoked := false
			target := func(*http.Request) (*http.Response, error) {
				targetInvoked = true
				return nil, nil
			}

			// When
			handler := flatten([]Filter{filter1, filter2}, target)

			// Then
			handler(nil)
			So(filter2Invoked, ShouldBeFalse)
			So(targetInvoked, ShouldBeFalse)
		})
	})
}
