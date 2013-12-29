package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestFronttier(t *testing.T) {
	Convey("Given a valid configuration, fronttier.json", t, func() {
		server, err := getServer("fronttier.json")

		Convey("Then I expect #getServer to retrieve a valid server instance", func() {
			So(err, ShouldBeNil)
			So(server, ShouldNotBeNil)
		})
	})
}
