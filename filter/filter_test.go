package filter

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/savaki/frontier/mock"
	"net/http"
)

type SampleFilter struct {
	value string
}

func (self *SampleFilter) Filter(w http.ResponseWriter, req *http.Request, target http.Handler) {
	key := "X-Key"
	value := req.Header.Get(key)
	if len(value) == 0 {
		req.Header.Set(key, self.value)
	} else {
		req.Header.Set(key, value+","+self.value)
	}

	target.ServeHTTP(w, req)
}

var _ = Describe("Filter", func() {
	Context("#Flatten", func() {
		It("should do something", func() {
			f1 := &SampleFilter{"a"}
			f2 := &SampleFilter{"b"}
			target := &mock.Handler{}

			request, _ := http.NewRequest("GET", "http://www.yahoo.com", nil)

			// When
			handler := Flatten([]HandlerFilter{f1, f2}, target)

			// Then
			handler.ServeHTTP(nil, request)
			Expect(request.Header.Get("X-Key")).To(Equal("a,b"), "expected headers to be applied in sequential order")
		})

		It("should return the handler if no filters were specified", func() {
			target := &mock.Handler{}

			// When
			handler := Flatten(nil, target)

			// Then
			Expect(handler).To(Equal(target), "expected our original handler back")
		})
	})
})
