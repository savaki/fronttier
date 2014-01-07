package main

import (
	"github.com/savaki/fronttier"
	"net/http"
)

func main() {
	router := fronttier.NewRouter()
	router.PathPrefix("/fantasy/").Proxy("http://sports.yahoo.com")
	router.PathPrefix("/").Proxy("http://www.cnn.com")

	http.ListenAndServe(":8080", router)
}
