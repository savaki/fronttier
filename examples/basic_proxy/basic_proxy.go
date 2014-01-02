package main

import (
	"github.com/savaki/fronttier"
)

func main() {
	builder := fronttier.Builder()
	builder.Paths("/fantasy/").Proxy().Url("http://sports.yahoo.com")
	builder.Paths("/").Proxy().Url("http://www.cnn.com")
	server, _ := builder.Build()

	server.ListenAndServe(":8080")
}
