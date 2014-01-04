package core

import (
	. "github.com/savaki/fronttier/matcher"
	"github.com/savaki/fronttier/mock"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"testing"
)

func TestRoute(t *testing.T) {
	var r *Route
	var handler http.Handler

	Convey("Given a Router", t, func() {
		handler = &mock.Handler{}
		r = &Route{
			matchers: []Matcher{Prefix("/sample")},
			handler:  handler,
		}

		Convey("When the matchers #Match the route", func() {
			req, _ := http.NewRequest("GET", "http://www.yahoo.com/sample", nil)
			matches := r.Matches(req)

			Convey("Then #Matches should return true", func() {
				So(matches, ShouldBeTrue)
			})
		})

		Convey("When something does something", func() {
			req, _ := http.NewRequest("GET", "http://www.yahoo.com/no-match", nil)
			matches := r.Matches(req)

			Convey("Then #Matches should return true", func() {
				So(matches, ShouldBeFalse)
			})
		})
	})
}
