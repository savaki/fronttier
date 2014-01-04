package fronttier

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestBuilder(t *testing.T) {
	Convey("When I call #Builder", t, func() {
		builder := Builder()

		Convey("Then I expect a *BuilderConfig instance", func() {
			So(builder, ShouldNotBeNil)
		})
	})
}
