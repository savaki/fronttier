package fronttier

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestProxy(t *testing.T) {
	Convey("When I call Proxy()", t, func() {
		Convey("Then I expect an handlerBuilder", func() {
			proxy := Proxy()

			So(proxy, ShouldNotBeNil)
		})
	})
}
