package mock

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestResponseWriter(t *testing.T) {
	Convey("#Header", t, func() {
		Convey("should capture headers", func() {
			w := &ResponseWriter{}

			// When
			w.Header().Set("Hello", "world")

			// Then
			So(w.Header().Get("Hello"), ShouldEqual, "world")
		})
	})

	Convey("#WriteHeader", t, func() {
		Convey("should capture status code", func() {
			w := &ResponseWriter{}
			statusCode := 201

			// When
			w.WriteHeader(statusCode)

			// Then
			So(w.StatusCode, ShouldEqual, statusCode)
		})
	})

	Convey("#Write", t, func() {
		Convey("should save the content", func() {
			w := &ResponseWriter{}
			content := []byte("hello world")

			// When
			w.Write(content)

			// Then
			So(string(w.Content), ShouldEqual, string(content))
		})
	})
}
