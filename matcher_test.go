package fronttier

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"testing"
)

func TestMatchers(t *testing.T) {
	var matcher *OrMatcher
	var req *http.Request

	Convey("Given an OrMatcher", t, func() {
		matcher = &OrMatcher{}
		req, _ = http.NewRequest("GET", "http://www.yahoo.com/sample", nil)

		Convey("When there are zero matchers", func() {
			Convey("Then #Matches should fail", func() {
				result := matcher.Matches(req)

				So(result, ShouldBeFalse)
			})
		})

		Convey("When at least one matcher matches", func() {
			matcher.Add(&PrefixMatcher{"/sam"})
			matcher.Add(&PrefixMatcher{"/boy"})
			matcher.Add(&PrefixMatcher{"/girl"})

			result := matcher.Matches(req)
			Convey("Then matcher should pass", func() {
				So(result, ShouldBeTrue)
			})
		})

		Convey("When I use the Or helper to construct the matcher", func() {
			orMatcher := Or(&PrefixMatcher{"/argle"}, &PrefixMatcher{"/bargle"}, &PrefixMatcher{"/sam"})

			Convey("Then I expect when at least one matcher matches, then matcher should pass", func() {
				result := orMatcher.Matches(req)

				So(result, ShouldBeTrue)
			})

			Convey("Then I expect when no matcher matches, then matcher should fail", func() {
				req, _ = http.NewRequest("GET", "http://www.funky.com/whoa!", nil)
				result := orMatcher.Matches(req)

				So(result, ShouldBeFalse)
			})
		})
	})
}
