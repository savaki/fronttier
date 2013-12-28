package proxy

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"testing"
)

func TestService(t *testing.T) {
	Convey("#authorized", t, func() {
		Convey("When NO requiredHeaders are defined", func() {
			handler := &proxyService{}
			request, _ := http.NewRequest("GET", "http://www.google.com", nil)

			// When
			result := handler.authorized(request)

			Convey("Then authorized should return true", func() {
				So(result, ShouldBeTrue)
			})
		})

		Convey("When ALL required headers are present", func() {
			required := "X-User-Id"
			handler := &proxyService{}
			handler.requiredHeaders = []string{required}
			request, _ := http.NewRequest("GET", "http://www.google.com", nil)
			request.Header.Set(required, "blah")

			// When
			result := handler.authorized(request)

			// Then
			Convey("Then #authorized should return true", func() {
				So(result, ShouldBeTrue)
			})
		})

		Convey("When some of the required headers are missing", func() {
			required1 := "X-Required-1"
			required2 := "X-Required-2"
			handler := &proxyService{}
			handler.requiredHeaders = []string{required1, required2}
			request, _ := http.NewRequest("GET", "http://www.google.com", nil)
			request.Header.Set(required1, "blah")

			// When
			result := handler.authorized(request)

			Convey("Then #authorized should return false", func() {
				So(result, ShouldBeFalse)
			})
		})
	})
}
