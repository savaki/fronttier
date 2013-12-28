package filter

import (
	"net/http"
)

type HandlerFilter interface {
	Filter(http.ResponseWriter, *http.Request, http.Handler)
}

type HandlerFilterAdapter struct {
	target http.Handler
	filter HandlerFilter
}

func (self *HandlerFilterAdapter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	self.filter.Filter(w, req, self.target)
}

func Flatten(filters []HandlerFilter, target http.Handler) http.Handler {
	if filters == nil || len(filters) == 0 {
		return target
	}

	for i := len(filters) - 1; i >= 0; i-- {
		safeTarget := target
		adapter := &HandlerFilterAdapter{
			target: safeTarget,
			filter: filters[i],
		}
		target = adapter
	}

	return target
}
