package main

import (
	. "github.com/savaki/fronttier"
)

func main() {
	builder := Builder()
	builder.Paths("/fantasy/").Handler(Proxy().Url("http://sports.yahoo.com"))
	builder.Paths("/").Handler(Proxy().Url("http://www.cnn.com"))
	server, _ := builder.Build()

	server.ListenAndServe(":8080")
}
