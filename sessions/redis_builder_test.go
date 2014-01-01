package sessions

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestRedisBuilder(t *testing.T) {
	var config *RedisConfig

	Convey("Given a redis builder", t, func() {
		config = Redis()

		Convey("When I don't specify any configurations", func() {
			Convey("Then I expect redis to default to localhost", func() {
				So(config.getHost(), ShouldEqual, ":6379")
			})

			Convey("Then I expect redis to default the namespace to sessions", func() {
				So(config.getNamespace(), ShouldEqual, "sessions")
			})

			Convey("Then I expect redis to default the maxIdle to 3", func() {
				So(config.getMaxIdle(), ShouldEqual, 3)
			})

			Convey("Then I expect a connFactory to be set", func() {
				So(config.getConnFactory(), ShouldNotBeNil)
			})
		})

		Convey("When I specify the Host", func() {
			host := "foo:6379"
			config.Host(host)

			Convey("Then I expect the specified host to be used", func() {
				So(config.getHost(), ShouldEqual, host)
			})
		})

		Convey("When I specify the Namespace", func() {
			namespace := "eek"
			config.Namespace(namespace)

			Convey("Then I expect the specified namespace to be used", func() {
				So(config.getNamespace(), ShouldEqual, namespace)
			})
		})

		Convey("When I specify MaxIdle", func() {
			maxIdle := 7
			config.MaxIdle(maxIdle)

			Convey("Then I expect the specified maxIdle to be used", func() {
				So(config.getMaxIdle(), ShouldEqual, maxIdle)
			})
		})
	})
}
