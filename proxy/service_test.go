package proxy

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
)

var _ = Describe("Service", func() {
	Context("#authorized", func() {
		It("should return true if no requiredHeaders defined", func() {
			handler := &proxyService{}
			request, _ := http.NewRequest("GET", "http://www.google.com", nil)

			// When
			result := handler.authorized(request)

			// Then
			Expect(result).To(Equal(true), "expected true when no required headers present")
		})

		It("should return true if all the required headers are present", func() {
			required := "X-User-Id"
			handler := &proxyService{}
			handler.requiredHeaders = []string{required}
			request, _ := http.NewRequest("GET", "http://www.google.com", nil)
			request.Header.Set(required, "blah")

			// When
			result := handler.authorized(request)

			// Then
			Expect(result).To(Equal(true), "expected true when no required headers present")
		})

		It("should return false if the required header is missing", func() {
			required := "X-Required"
			handler := &proxyService{}
			handler.requiredHeaders = []string{required}
			request, _ := http.NewRequest("GET", "http://www.google.com", nil)
			request.Header.Set(required, "blah")

			// When
			result := handler.authorized(request)

			// Then
			Expect(result).To(Equal(true), "expected true when our header is present")
		})

		It("should return false if some of the required headers are missing", func() {
			required1 := "X-Required-1"
			required2 := "X-Required-2"
			handler := &proxyService{}
			handler.requiredHeaders = []string{required1, required2}
			request, _ := http.NewRequest("GET", "http://www.google.com", nil)
			request.Header.Set(required1, "blah")

			// When
			result := handler.authorized(request)

			// Then
			Expect(result).To(Equal(false), "expected false when all our required headers are not present")
		})
	})
})
