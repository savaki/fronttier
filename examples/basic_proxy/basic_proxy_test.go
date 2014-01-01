package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestExample(t *testing.T) {
	Convey("Given this test runs", t, func() {
		Convey("Then we know the example compiles", func() {
			So(true, ShouldBeTrue)
		})
	})
}
