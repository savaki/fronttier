frontier
========

[![GoDoc](https://godoc.org/github.com/savaki/fronttier?status.png)](https://godoc.org/github.com/savaki/fronttier) [![Build Status](https://travis-ci.org/savaki/fronttier.png)](https://travis-ci.org/savaki/fronttier) [![Coverage Status](https://coveralls.io/repos/savaki/fronttier/badge.png)](https://coveralls.io/r/savaki/fronttier)

A routing reverse-proxy written in Go.

## Features

Fronttier is designed to be a front facing reverse-proxy.  Requests can be configured to route to different services e.g. google.com/mail, google.com/analytics, google.com/adwords, etc.

* routes requests to different services
* can be optionally configured to also manage user sessions 

### Service Routing

Let's start with a simple example.  

#### Example - Simple Routing

Suppose our site consists of three services as shown in the following diagram:

![image](Overview.png)

We could configure Fronttier as follows:

```
package main

import (
  . "github.com/savaki/fronttier"
)

func main() {
  builder := Builder()

  builder.AuthConfig().ReservedHeaders("X-User-Id", "X-Name", "X-Email")

  builder.Paths("/x").Handler(Proxy().Url("http://x-service"))
  builder.Paths("/y").Handler(Proxy().Url("http://y-service"))
  builder.
    Paths("/login").
    Handler(Proxy().Url("http://login-service"))

  server, _ := builder.Build()
  server.ListenAndServe(":8080")
}
```

### Session Management 

In addition to routing between sites, fronttier can also provide session management to those services.  Fronttier can be configured to do this by calling #SessionFactory on the path.  

#### Configuring:

1. Mark one or more routes with #SessionFactory
2. Optionally define one or more reserved headers.

#### Creating Sessions:

If a route marked with #SessionFactory returns a reserved header:

1. create a new session
2. place all the reserved headers and their values into the session
3. return to the user a cookie that identifies the session

![image](Login.png)

#### Handling Requests:

Once a session has been created, if a request is received that contains a valid session cookie:

1. the reserved headers will be retrieved from the service
2. and added to the request when sent the underlying service

![image](Request.png)

#### Logging Out:

To logout, any service may send the X-Logout header:

1. delete the session from the session store
2. clear the cookie

![image](Logout.png)

#### Protecting Against Forgery:

To defend against forgery, any reserved header received from the browser will be removed.  In the future, fronttier may include a signed header similar to how Amazon handles security.

![image](Forgery.png)

## Sample Code

We can modify our previous example to this:

```
package main

import (
  . "github.com/savaki/fronttier"
)

func main() {
  builder := Builder()

  builder.AuthConfig().ReservedHeaders("X-User-Id", "X-Name", "X-Email")

  builder.Paths("/x").Handler(Proxy().Url("http://x-service"))
  builder.Paths("/y").Handler(Proxy().Url("http://y-service"))
  builder.
    Paths("/login").
    SessionFactory(). // mark this route as capable of creating a session
    Handler(Proxy().Url("http://login-service"))

  server, _ := builder.Build()
  server.ListenAndServe(":8080")
}
```


# Testing

Fronttier uses the very excellent [goconvey](https://github.com/smartystreets/goconvey) framework for testing.  To see all the tests execute in a browser:

```
go get github.com/smartystreets/goconvey
goconvey &
open http://localhost:8080
```
