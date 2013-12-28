package mock

import (
	"errors"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"testing"
)

func TestRoundTripper(t *testing.T) {
	Convey("#RoundTrip", t, func() {
		Convey("should capture the request and return the response", func() {
			request, _ := http.NewRequest("GET", "http://www.sample.com/abc", nil)
			tripper := &RoundTripper{
				Response: &http.Response{StatusCode: 201},
				Err:      errors.New("blah"),
			}

			// When
			response, err := tripper.RoundTrip(request)

			// Then
			So(tripper.Request, ShouldEqual, request)
			So(response, ShouldEqual, tripper.Response)
			So(err, ShouldEqual, tripper.Err)
		})
	})
}
