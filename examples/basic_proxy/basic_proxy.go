package main

import (
	"github.com/savaki/fronttier"
	"github.com/savaki/fronttier/proxy"
)

func main() {
	yahoo, _ := proxy.Builder().
		Url("http://sports.yahoo.com").
		Build()

	cnn, _ := proxy.Builder().
		Url("http://www.cnn.com").
		Build()

	builder := frontier.Builder()
	builder.Path("/fantasy/").Handler(yahoo)
	builder.Path("/").Handler(cnn)
	server, _ := builder.Build()

	server.ListenAndServe(":8080")
}
