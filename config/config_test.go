package config

import (
	"encoding/json"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestConfig(t *testing.T) {
	var cfg *Config

	Convey("Given a valid json file", t, func() {
		data := []byte(`{"routes":[
			{
				"paths": ["/"],
				"proxy": "http://www.google.com"
			}
		]}`)

		Convey("Then I expect to read the data into Config", func() {
			cfg = &Config{}
			err := json.Unmarshal(data, cfg)

			So(err, ShouldBeNil)
			So(cfg, ShouldNotBeNil)
			So(len(cfg.Routes), ShouldEqual, 1)

			route := cfg.Routes[0]
			So(len(route.Paths), ShouldEqual, 1)
			So(route.Paths[0], ShouldEqual, "/")
		})
	})

	Convey("#Validate", t, func() {
		cfg = &Config{}

		Convey("Given a config that has zero routes", func() {
			data := []byte(`{"routes":[]}`)
			json.Unmarshal(data, cfg)
			err := cfg.Validate()

			Convey("Then I expect #Validate to return an error", func() {
				So(err, ShouldNotBeNil)
			})
		})

		Convey("Given a valid config", func() {
			data := []byte(`{"routes":[
				{
					"paths": ["/"],
					"proxy": "http://www.google.com"
				}
			]}`)
			json.Unmarshal(data, cfg)
			err := cfg.Validate()

			Convey("Then I expect #Validate to return nil", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}
