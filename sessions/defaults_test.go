package sessions

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestDefaultIdFactory(t *testing.T) {
	Convey("Given a defaultIdFactory", t, func() {
		Convey("When I invoke it", func() {
			result := defaultIdFactory()

			Convey("Then I expect it to generate a new id", func() {
				So(len(result), ShouldBeGreaterThan, 0)
			})
		})
	})
}
