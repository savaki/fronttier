package proxy

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/savaki/fronttier/mock"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"testing"
)

func TestBuilder(t *testing.T) {
	Convey("#AddFilter", t, func() {

	})
}

var _ = Describe("Builder", func() {
	Context("#AddFilter", func() {
		It("add filters to the handler", func() {
			var filterInvoked *bool = new(bool) // check to see filter was invoked
			*filterInvoked = false

			tripper := &mock.RoundTripper{}
			filter := func(req *http.Request, next Handler) (*http.Response, error) {
				*filterInvoked = true
				return nil, errors.New("blah")
			}

			handler, err := Builder().Url("http://www.cnn.com").AddFilter(filter).RoundTripper(tripper).Build()
			Expect(err).To(BeNil(), "expected to construct a handler")

			// When
			request, _ := http.NewRequest("GET", "http://www.google.com", nil)
			handler.ServeHTTP(nil, request)

			// Then
			Expect(*filterInvoked).To(BeTrue(), "expected the filter to have been invoked")
			Expect(tripper.Request).To(BeNil(), "expected request to not have made it to the RoundTripper")
		})
	})

	Context("#Build", func() {
		It("returns error if Url not defined", func() {
			_, err := Builder().Build()

			// Then
			Expect(err).ToNot(BeNil(), "expected an error to be thrown since Url wasn't defined")
		})

		It("return service if Url is defined", func() {
			service, err := Builder().Url("http://www.google.com").Build()

			Expect(err).To(BeNil(), "expected no errors")
			Expect(service).ToNot(BeNil(), "expected a service to be defined")
		})

		It("sets the Host header for requests that are passed through", func() {
			tripper := mock.RoundTripper{}
			service, _ := Builder().Url("http://www.google.com").RoundTripper(&tripper).Build()
			request, _ := http.NewRequest("GET", "http://www.yahoo.com/sample", nil)

			// When
			service.ServeHTTP(nil, request)

			// Then
			Expect(tripper.Request.URL.Host).To(Equal("www.google.com"), "expected to connect to proxied host")
			Expect(tripper.Request.Header.Get("Host")).To(Equal("www.google.com"), "expected host header to be set")
		})

		It("should report errors thrown earlier in the process", func() {
			_, err := Builder().Url(":::this is junk").Build()

			// Then
			Expect(err).ToNot(BeNil(), "expected an error to be returned")
		})

		It("the created Handler should return prematurely if the RoundTripper fails", func() {
			tripper := &mock.RoundTripper{
				Request:  nil,
				Response: nil,
				Err:      errors.New("boom!"),
			}
			service, _ := Builder().RoundTripper(tripper).Url("http://www.google.com").Build()
			request, _ := http.NewRequest("GET", "http://www.cnn.com", nil)

			// When
			service.ServeHTTP(nil, request)

			// Then
			// the response writer should not be written to
		})
	})

	Context("#NotAuthorizedHandler", func() {
		It("should invoke the default not authorized handler when required headers not provided", func() {
			notAuthorized := &defaultNotAuthorizedHandler{}
			service, err := Builder().
				Url("http://www.cnn.com").
				RequiredHeaders("X-Foo").
				NotAuthorizedHandler(notAuthorized).
				Build()
			Expect(err).To(BeNil(), "expected no errors")

			// When
			req, _ := http.NewRequest("GET", "http://www.cnn.com", nil)
			w := &mock.ResponseWriter{}
			service.ServeHTTP(w, req)

			// Then
			Expect(w.StatusCode).To(Equal(401), "expected a 401 status code")
		})
	})

	Context("#RequiredHeaders", func() {
		It("should optionally allow authentication to be required by passing a known header", func() {
			service, err := Builder().Url("http://www.cnn.com/").RequiredHeaders("X-Foo").Build()

			// Then
			Expect(err).To(BeNil(), "expected no errors")
			Expect(service).ToNot(BeNil(), "expected a service back")
		})

		It("should assign required headers in the service", func() {
			handler, err := Builder().Url("http://www.cnn.com/").RequiredHeaders("X-Foo").Build()

			// Then
			s := handler.(*proxyService)
			Expect(err).To(BeNil(), "expected no errors")
			Expect(s.requiredHeaders).To(Equal([]string{"X-Foo"}), "expected our required header to have been set")
		})

		It("should assign a default NotAuthorizedHandler if one wasn't already defined", func() {
			service, err := Builder().
				Url("http://www.cnn.com").
				RequiredHeaders("X-Foo").
				Build()
			Expect(err).To(BeNil(), "expected no errors")

			// When
			req, _ := http.NewRequest("GET", "http://www.cnn.com", nil)
			w := &mock.ResponseWriter{}
			service.ServeHTTP(w, req)

			// Then
			Expect(w.StatusCode).To(Equal(401), "expected a 401 status code")
		})
	})
})
