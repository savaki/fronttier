package fronttier

import (
	"net/http"
)

type FilterFunc func(http.ResponseWriter, *http.Request, http.HandlerFunc)

type FilterWithNameFunc func(http.ResponseWriter, *http.Request, http.HandlerFunc, string)

func flatten(filters []FilterFunc, target http.HandlerFunc) http.HandlerFunc {
	if filters == nil || len(filters) == 0 {
		return target
	}

	for i := len(filters) - 1; i >= 0; i-- {
		filter := filters[i]
		safeTarget := target
		target = func(w http.ResponseWriter, req *http.Request) {
			filter(w, req, safeTarget)
		}
	}

	return target
}
