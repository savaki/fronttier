package proxy

import (
	"net/http"
)

type Handler func(req *http.Request) (*http.Response, error)

type Filter func(req *http.Request, target Handler) (*http.Response, error)

func flatten(filters []Filter, target Handler) Handler {
	if filters == nil || len(filters) == 0 {
		return target
	}

	for i := len(filters) - 1; i >= 0; i-- {
		filter := filters[i]
		safeTarget := target
		target = func(req *http.Request) (*http.Response, error) {
			return filter(req, safeTarget)
		}
	}

	return target
}
