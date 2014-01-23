package sessions

import (
	"errors"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"testing"
)

type MockSigner struct {
	InSign            *http.Request
	InSignRouteName   string
	InVerify          *http.Request
	InVerifyRouteName string
	OutVerify         bool
}

func (self *MockSigner) Sign(routeName string, req *http.Request) {
	self.InSignRouteName = routeName
	self.InSign = req
}

func (self *MockSigner) Verify(routeName string, req *http.Request) bool {
	self.InVerifyRouteName = routeName
	self.InVerify = req
	return self.OutVerify
}

func TestBuilder(t *testing.T) {
	var builder *BuilderConfig
	var signer *MockSigner
	var filter *AuthFilter
	var err error

	Convey("Given a builder", t, func() {
		builder = Builder().ReservedHeaders("X-Header")

		Convey("When the builder has an internal error", func() {
			builder.err = errors.New("I have an internal error")

			Convey("Then I expect #BuildAuthFilter() to fail", func() {
				_, err = builder.BuildAuthFilter()

				So(err, ShouldNotBeNil)
			})
		})

		Convey("When cookieName is assigned", func() {
			cookieName := "sample123"
			builder = builder.CookieName(cookieName)

			Convey("Then authFilter.template.Name should be cookieName", func() {
				filter, err = builder.BuildAuthFilter()

				So(err, ShouldBeNil)
				So(filter, ShouldNotBeNil)
				So(filter.template.Name, ShouldEqual, cookieName)
			})
		})

		Convey("When cookieName is NOT assigned", func() {
			Convey("Then authFilter.template.Name should default to 'id'", func() {
				filter, err = builder.BuildAuthFilter()

				So(err, ShouldBeNil)
				So(filter, ShouldNotBeNil)
				So(filter.template.Name, ShouldEqual, "id")
			})
		})

		Convey("When the signerFunc is assigned", func() {
			signer = &MockSigner{}
			filter, err = builder.SignerFunc(signer.Sign).BuildAuthFilter()

			Convey("Then authFilter.sign should be set", func() {
				So(err, ShouldBeNil)
				So(filter, ShouldNotBeNil)
				So(filter.sign, ShouldNotEqual, nil)
			})
		})

		Convey("When the #Signer is assigned", func() {
			signer = &MockSigner{}
			filter, err = builder.Signer(signer).BuildAuthFilter()

			Convey("Then authFilter.verify should be set", func() {
				So(err, ShouldBeNil)
				So(filter, ShouldNotBeNil)
				So(filter.sign, ShouldNotEqual, nil)
			})
		})

		Convey("When the #Signer is NOT assigned", func() {
			filter, err = builder.BuildAuthFilter()

			Convey("Then there should be no default value", func() {
				So(err, ShouldBeNil)
				So(filter, ShouldNotBeNil)
				So(filter.sign, ShouldEqual, nil)
			})
		})

		Convey("When the #Verifier is assigned", func() {
			signer = &MockSigner{}
			filter, err = builder.Verifier(signer).BuildAuthFilter()

			Convey("Then authFilter.verify should be set", func() {
				So(err, ShouldBeNil)
				So(filter, ShouldNotBeNil)
				So(filter.verify, ShouldNotEqual, nil)
			})
		})

		Convey("When the #VerifierFunc is assigned", func() {
			signer = &MockSigner{}
			filter, err = builder.VerifierFunc(signer.Verify).BuildAuthFilter()

			Convey("Then authFilter.verify should be set", func() {
				So(err, ShouldBeNil)
				So(filter, ShouldNotBeNil)
				So(filter.verify, ShouldNotEqual, nil)
			})
		})

		Convey("When the #Verifier is NOT assigned", func() {
			filter, err = builder.BuildAuthFilter()

			Convey("Then authFilter.verify should be set", func() {
				So(err, ShouldBeNil)
				So(filter, ShouldNotBeNil)
				So(filter.verify, ShouldEqual, nil)
			})
		})

		Convey("When cookieDomain is assigned", func() {
			cookieDomain := "blah.com"
			builder = builder.CookieDomain(cookieDomain)

			Convey("Then authFilter.template.Domain should be cookieDomain", func() {
				filter, err = builder.BuildAuthFilter()

				So(err, ShouldBeNil)
				So(filter, ShouldNotBeNil)
				So(filter.template.Domain, ShouldEqual, cookieDomain)
			})
		})

		Convey("When cookiePath is assigned", func() {
			cookiePath := "/blah"
			builder = builder.CookiePath(cookiePath)

			Convey("Then authFilter.template.Path should be cookiePath", func() {
				filter, err = builder.BuildAuthFilter()

				So(err, ShouldBeNil)
				So(filter, ShouldNotBeNil)
				So(filter.template.Path, ShouldEqual, cookiePath)
			})
		})

		Convey("When cookiePath is NOT assigned", func() {
			Convey("Then authFilter.template.Path should default to /", func() {
				filter, err = builder.BuildAuthFilter()

				So(err, ShouldBeNil)
				So(filter, ShouldNotBeNil)
				So(filter.template.Path, ShouldEqual, "/")
			})
		})

		Convey("When idFactory is assigned", func() {
			id := "hello world"
			idFactory := func() string {
				return id
			}
			builder = builder.IdFactory(idFactory)

			Convey("Then authFilter.idFactory should be idFactory", func() {
				filter, err = builder.BuildAuthFilter()

				So(err, ShouldBeNil)
				So(filter, ShouldNotBeNil)
				So(filter.idFactory(), ShouldEqual, id)
			})
		})

		Convey("When idFactory is NOT assigned", func() {
			Convey("Then authFilter.idFactory should use the default uuid factory", func() {
				filter, err = builder.BuildAuthFilter()

				So(err, ShouldBeNil)
				So(filter, ShouldNotBeNil)
				So(filter.idFactory, ShouldEqual, defaultIdFactory)
			})
		})

		Convey("When reservedHeaders is assigned", func() {
			header := "X-Uid"
			builder = builder.ReservedHeaders(header)

			Convey("Then authFilter.reservedHeaders should be reservedHeaders", func() {
				filter, err = builder.BuildAuthFilter()

				So(err, ShouldBeNil)
				So(filter, ShouldNotBeNil)
				So(filter.reservedHeaders, ShouldContain, header)
			})
		})

		Convey("When reservedHeaders is NOT assigned", func() {
			Convey("Then #BuildAuthFilter should fail", func() {
				_, err = Builder().BuildAuthFilter()

				So(err, ShouldNotBeNil)
			})
		})

		Convey("When logoutHeader is assigned", func() {
			header := "X-Blah"
			builder = builder.LogoutHeader(header)

			Convey("Then authFilter.logoutHeader should be logoutHeader", func() {
				filter, err := builder.BuildAuthFilter()

				So(err, ShouldBeNil)
				So(filter, ShouldNotBeNil)
				So(filter.logoutHeader, ShouldEqual, header)
			})
		})

		Convey("When logoutHeader is NOT assigned", func() {
			Convey("Then logoutHeader should default to X-Logout", func() {
				filter, err = builder.BuildAuthFilter()

				So(err, ShouldBeNil)
				So(filter, ShouldNotBeNil)
				So(filter.logoutHeader, ShouldEqual, "X-Logout")
			})
		})

		Convey("When sessionStore is assigned", func() {
			sessionStore := Memory()
			builder = builder.SessionStore(sessionStore)

			Convey("Then authFilter.reservedHeaders should be reservedHeaders", func() {
				filter, err = builder.BuildAuthFilter()

				So(err, ShouldBeNil)
				So(filter, ShouldNotBeNil)
				So(filter.sessionStore, ShouldEqual, sessionStore)
			})
		})

		Convey("When sessionStore is NOT assigned", func() {
			Convey("Then the default session store should be assigned", func() {
				filter, err := builder.BuildAuthFilter()

				So(err, ShouldBeNil)
				So(filter.sessionStore, ShouldEqual, defaultSessionStore)
			})
		})

		Convey("When timeoutMinutes is assigned", func() {
			timeoutMinutes := 7
			builder = builder.TimeoutMinutes(timeoutMinutes)

			Convey("Then authFilter.timeoutMinutes should be timeoutMinutes", func() {
				filter, err = builder.BuildAuthFilter()

				So(err, ShouldBeNil)
				So(filter, ShouldNotBeNil)
				So(filter.template.MaxAge, ShouldEqual, timeoutMinutes*60)
			})
		})

		Convey("When timeoutMinutes is NOT assigned", func() {
			Convey("Then timeoutMinutes should default to 15", func() {
				filter, err = builder.BuildAuthFilter()

				So(err, ShouldBeNil)
				So(filter.template.MaxAge, ShouldEqual, 15*60)
			})
		})
	})
}
