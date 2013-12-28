package mock

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ResponseWriter", func() {
	Context("#Header", func() {
		It("should capture headers", func() {
			w := &ResponseWriter{}

			// When
			w.Header().Set("Hello", "world")

			// Then
			Expect(w.Header().Get("Hello")).To(Equal("world"), "expected header to be set")
		})
	})

	Context("#WriteHeader", func() {
		It("should capture status code", func() {
			w := &ResponseWriter{}
			statusCode := 201

			// When
			w.WriteHeader(statusCode)

			// Then
			Expect(w.StatusCode).To(Equal(statusCode), "expected status code to be set")
		})
	})

	Context("#Write", func() {
		It("should save the content", func() {
			w := &ResponseWriter{}
			content := []byte("hello world")

			// When
			w.Write(content)

			// Then
			Expect(w.Content).To(Equal(content), "expected content to be captured")
		})
	})
})
