package mock

import (
	"bytes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
)

var _ = Describe("Handler", func() {
	Context("#ServeHTTP", func() {
		It("captures the request path", func() {
			request, _ := http.NewRequest("GET", "http://www.google.com/mail/", nil)

			// When
			handler := &Handler{}
			handler.ServeHTTP(nil, request)

			// Then
			Expect(handler.InPath).To(Equal("/mail/"), "expected path to be captured")
		})

		It("captures the request content", func() {
			content := "hello world"
			request, _ := http.NewRequest("GET", "http://www.google.com/mail/", bytes.NewReader([]byte(content)))

			// When
			handler := &Handler{}
			handler.ServeHTTP(nil, request)

			// Then
			Expect(string(handler.InContent)).To(Equal(content), "expected to capture request content")
		})

		It("expected headers and values to be saved by handler", func() {
			request, _ := http.NewRequest("GET", "http://www.google.com/mail/", nil)
			request.Header.Set("hello", "world")

			// When
			handler := &Handler{}
			handler.ServeHTTP(nil, request)

			// Then
			Expect(handler.InMethod).To(Equal("GET"), "expected method to be captured")
			Expect(handler.InHeader["Hello"]).To(Equal([]string{"world"}), "expected headers to be captured")
		})

		It("write the specified status code", func() {
			handler := &Handler{OutStatusCode: 201}
			w := &ResponseWriter{}

			// When
			handler.ServeHTTP(w, nil)

			// Then
			Expect(w.StatusCode).To(Equal(handler.OutStatusCode), "expected status code to be returned")
		})

		It("defaults status code to 200", func() {
			handler := &Handler{}
			w := &ResponseWriter{}

			// When
			handler.ServeHTTP(w, nil)

			// Then
			Expect(w.StatusCode).To(Equal(200), "expected 200 status code as default")
		})

		It("returns the out headers", func() {
			header := "X-Hello"
			handler := &Handler{OutHeader: map[string]string{header: "world"}}
			w := &ResponseWriter{}

			// When
			handler.ServeHTTP(w, nil)

			// Then
			Expect(w.Header().Get(header)).To(Equal(handler.OutHeader[header]), "expected the out header to have been set")
		})

		It("returns the out content", func() {
			content := []byte("hello world")
			handler := &Handler{OutContent: content}
			w := &ResponseWriter{}

			// When
			handler.ServeHTTP(w, nil)

			// Then
			Expect(w.Content).To(Equal(handler.OutContent), "expected the out content to have been returned")
		})
	})
})
