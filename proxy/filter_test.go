package proxy

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
)

var _ = Describe("Filter", func() {
	Context("#flatten", func() {
		It("should return the target if no filtering needs to happen", func() {
			gripe := errors.New("blah")
			var filters []Filter
			var target = func(*http.Request) (*http.Response, error) {
				return nil, gripe
			}

			// When
			result := flatten(filters, target)

			// Then
			response, err := result(nil)
			Expect(response).To(BeNil(), "expected no response")
			Expect(err).To(Equal(gripe), "expected our error to have been thrown")
		})

		It("should process filters before the target", func() {
			gripe := errors.New("blah")
			filter := func(req *http.Request, next Handler) (*http.Response, error) {
				return nil, gripe
			}
			filters := []Filter{filter}
			target := func(*http.Request) (*http.Response, error) {
				Fail("target should not have been called")
				return nil, nil
			}

			// When
			handler := flatten(filters, target)
			handler(nil)
		})

		It("should process the FIRST filter provided FIRST", func() {
			gripe := errors.New("blah")
			filter1 := func(req *http.Request, next Handler) (*http.Response, error) {
				return nil, gripe
			}
			filter2 := func(req *http.Request, next Handler) (*http.Response, error) {
				Fail("second filter should not have been called")
				return nil, nil
			}
			target := func(*http.Request) (*http.Response, error) {
				Fail("target handler should not have been called")
				return nil, nil
			}

			// When
			handler := flatten([]Filter{filter1, filter2}, target)

			// Then
			handler(nil)
		})
	})
})
