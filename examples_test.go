package fronttier

func ExampleBuilder() {
	builder := Builder()
	builder.Paths("/fantasy/").Handler(Proxy().Url("http://sports.yahoo.com"))
	builder.Paths("/").Handler(Proxy().Url("http://www.cnn.com"))
	server, _ := builder.Build()
	server.ListenAndServe(":8080")
}

// func ExampleFoo() {
// 	route1 := Route("/fantasy/").Method("POST").Proxy("http://sports.yahoo.com")
// 	route2 := Route("/").Method("POST").Proxy("http://www.cnn.com")
// 	session := Session().CookieName("sample").Provider(REDIS).Host("foo.bar")

// 	builder := Builder().
// 		Handle("/fantasy", Proxy().Url("http://sports.yahoo.com")).
// 		HandleFunc("/fantasy", func(w http.ResponseWriter, req *http.Request) {}).
// 		HandleRoute(route1).
// 		HandleRoute(route2).
// 		HandleSession(session).
// 		Build()
// }
