package fronttier

import (
	"bytes"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestConfig(t *testing.T) {
	Convey("#Load", t, func() {
		Convey("When I load a badly formed json file", func() {
			badlyFormed := []byte(`{`)
			builder := Load(bytes.NewReader(badlyFormed))

			Convey("Then I expect the builder to have an error", func() {
				So(builder.err, ShouldNotBeNil)
			})
		})

		Convey("When I load a valid json file that has an incomplete route configuration", func() {
			incomplete := []byte(`{"routes":[
				{
					"paths":["/"]
				}
			]}`)
			builder := Load(bytes.NewReader(incomplete))

			Convey("Then I expect the builder to have an error", func() {
				So(builder.err, ShouldNotBeNil)
			})
		})
	})

	Convey("#LoadFile", t, func() {
		Convey("When I attempt to load a file that doesn't exist", func() {
			builder := LoadFile("does-not-exist")

			Convey("Then I expect the builder to have an error", func() {
				So(builder.err, ShouldNotBeNil)
			})
		})

		Convey("When I load a valid configuration from file, sample.json", func() {
			builder := LoadFile("config/sample.json")

			Convey("Then I expect a valid builder", func() {
				So(builder, ShouldNotBeNil)
				So(builder.err, ShouldBeNil)
			})

			Convey("And I expect the builder to have 1 route config", func() {
				So(len(builder.routeConfigs), ShouldEqual, 1)
			})
		})
	})
}
