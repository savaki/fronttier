package mock

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
)

var _ = Describe("RoundTripper", func() {
	Context("#RoundTrip", func() {
		It("should capture the request and return the response", func() {
			request, _ := http.NewRequest("GET", "http://www.sample.com/abc", nil)
			tripper := &RoundTripper{
				Response: &http.Response{StatusCode: 201},
				Err:      errors.New("blah"),
			}

			// When
			response, err := tripper.RoundTrip(request)

			// Then
			Expect(tripper.Request).To(Equal(request), "expected request to be captured")
			Expect(response).To(Equal(tripper.Response), "expected response to be returned")
			Expect(err).To(Equal(tripper.Err), "expected err to be returned")
		})
	})
})
