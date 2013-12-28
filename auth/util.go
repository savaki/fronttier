package auth

import (
	"net/http"
)

func newCookie(template *http.Cookie, value string) *http.Cookie {
	return &http.Cookie{
		Name:   template.Name,
		Value:  value,
		Path:   template.Path,
		Domain: template.Domain,
		MaxAge: template.MaxAge,
	}
}

func copyNonReservedHeaders(reservedHeaders []string, source http.ResponseWriter, target http.ResponseWriter) {
Loop:
	for key, values := range source.Header() {
		// strip reserved headers from response
		for _, reservedHeader := range reservedHeaders {
			if key == reservedHeader {
				continue Loop
			}
		}

		for _, value := range values {
			target.Header().Add(key, value)
		}
	}
}
