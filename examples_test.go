package fronttier

import (
	"net/http"
)

func ExampleNewRouter() {
	router := NewRouter()

	router.Sessions().ReservedHeaders("X-User-Id", "X-Name", "X-Email")

	router.PathPrefix("/x").Proxy("http://x-service")
	router.PathPrefix("/y").Proxy("http://y-service")

	router.NewRoute().
		PathPrefix("/login").
		SessionFactory().
		Proxy("http://login-service")

	http.ListenAndServe(":8080", router)
}
