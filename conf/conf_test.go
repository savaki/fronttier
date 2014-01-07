package conf

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestLoadFile(t *testing.T) {
	var config *Config
	var err error

	Convey("Given a config file", t, func() {
		filename := "config.json"

		Convey("When I call #LoadFile to load the config", func() {
			config, err = LoadFile(filename)

			Convey("Then I expect no err to be thrown", func() {
				So(err, ShouldBeNil)
			})

			Convey("And I expect a valid *Config back", func() {
				So(config, ShouldNotBeNil)
			})

			Convey("And I expect 2 routes to be defined", func() {
				So(len(config.Routes), ShouldEqual, 2)
			})

			Convey("And I the first route's Paths to be /fantasy/", func() {
				So(config.Routes[0].Paths, ShouldEqual, "/fantasy/")
			})

			Convey("And I the first route to be a #SessionFactory", func() {
				So(config.Routes[0].SessionFactory, ShouldBeTrue)
			})

			Convey("And I the cookie name to be my-cookie", func() {
				So(config.Sessions.CookieName, ShouldEqual, "my-cookie")
			})
		})

		Convey("When I call #LoadFile to load a file that doesn't exist", func() {
			config, err = LoadFile("does-not-exist")

			Convey("Then I expect an error", func() {
				So(err, ShouldNotBeNil)
				So(config, ShouldBeNil)
			})
		})
	})
}
