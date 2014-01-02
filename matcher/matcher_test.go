package matcher

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"testing"
)

func TestMatchers(t *testing.T) {
	var matcher Matcher
	var req *http.Request

	Convey("Given an orMatcher", t, func() {
		req, _ = http.NewRequest("GET", "http://www.yahoo.com/sample", nil)

		Convey("When there are zero matchers", func() {
			matcher = Or()

			Convey("Then #Matches should fail", func() {
				result := matcher.Matches(req)

				So(result, ShouldBeFalse)
			})
		})

		Convey("When there is exactly ONE matcher", func() {
			prefix := Prefix("/hello")
			matcher = Or(prefix)

			Convey("Then I expect Or to return the one matcher", func() {
				So(matcher, ShouldResemble, prefix)
			})
		})

		Convey("When there are multiple matchers", func() {
			matcher = Or(Prefix("/sam"), Prefix("/boy"), Prefix("/girl"))

			result := matcher.Matches(req)
			Convey("Then matcher should pass", func() {
				So(result, ShouldBeTrue)
			})
		})

		Convey("When I use the Or helper to construct the matcher", func() {
			matcher := Or(Prefix("/argle"), Prefix("/bargle"), Prefix("/sam"))

			Convey("Then I expect when at least one matcher matches, then matcher should pass", func() {
				result := matcher.Matches(req)

				So(result, ShouldBeTrue)
			})

			Convey("Then I expect when no matcher matches, then matcher should fail", func() {
				req, _ = http.NewRequest("GET", "http://www.funky.com/whoa!", nil)
				result := matcher.Matches(req)

				So(result, ShouldBeFalse)
			})
		})
	})
}
