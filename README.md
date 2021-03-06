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

![image](docs/Overview.png)

We could configure Fronttier as follows:

```
package main

import (
  . "github.com/savaki/fronttier"
)

func main() {
  router := fronttier.NewRouter()

  router.PathPrefix("/x").Proxy("http://x-service")
  router.PathPrefix("/y").Proxy("http://y-service")

  router.NewRoute().
    PathPrefix("/login").
    Proxy("http://login-service")

  http.ListenAndServe(":8080", router)
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

![image](docs/Login.png)

#### Handling Requests:

Once a session has been created, if a request is received that contains a valid session cookie:

1. the reserved headers will be retrieved from the service
2. and added to the request when sent the underlying service

![image](docs/Request.png)

#### Logging Out:

To logout, any service may send the X-Logout header:

1. delete the session from the session store
2. clear the cookie

![image](docs/Logout.png)

#### Protecting Against Forgery:

To defend against forgery, any reserved header received from the browser will be removed.  In the future, fronttier may include a signed header similar to how Amazon handles security.

![image](docs/Forgery.png)

## Sample Code

We can modify our previous example to this:

```
package main

import (
  . "github.com/savaki/fronttier"
)

func main() {
  router := fronttier.NewRouter()

  router.Sessions().ReservedHeaders("X-User-Id", "X-Name", "X-Email")

  router.PathPrefix("/x").Proxy("http://x-service")
  router.PathPrefix("/y").Proxy("http://y-service")

  router.NewRoute().
    PathPrefix("/login").
    SessionFactory().
    Proxy("http://login-service")

  http.ListenAndServe(":8080", router)
}
```


# Testing

Fronttier uses the very excellent [goconvey](https://github.com/smartystreets/goconvey) framework for testing.  To see all the tests execute in a browser:

```
go get github.com/smartystreets/goconvey
goconvey &
open http://localhost:8080
```
