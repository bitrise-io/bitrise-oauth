[![Build Status](https://app.bitrise.io/app/e6a7166eda823c72/status.svg?token=LACL0_krbTkiMlmi4kBLNA&branch=master)](https://app.bitrise.io/app/e6a7166eda823c72)

# Bitrise OAuth library for Go
This is a very thin package over Go's standard [OAuth2 library](https://github.com/golang/oauth2), extending its functionality via introducing an additional layer that handles the initialization and communication with our current authorization provider [Keycloak](https://github.com/keycloak/keycloak).

This package provides both *client-side* and *server-side* (covering all of our current use-cases) wrappers. In this document you may find useful information about the APIs, the custom configuration options and the usage as well.

## Client
lorem ipsum intro




## Server
The server-side validation logic is located at the `service` package. You can use the `Validator` to via several different methods in order to validate any request. The supported use-cases are the following:
- **Handler Function** with:
	- Go's default HTTP multiplexer
	- Gorilla's HTTP router called [**gorilla/mux**](https://github.com/gorilla/mux)
- **Middleware** with:
	- Go's default HTTP multiplexer
	- Gorilla's HTTP router called [**gorilla/mux**](https://github.com/gorilla/mux)
- **Middleware Function** and **Handler Function** with Labstack's router called [**echo**](https://github.com/labstack/echo)


### API
for service package


### Options
The package offers wide configurability using Options. You can easily override any parameter passing the desired Option(s) as a constructor parameter. Not only the `Validator` itself have Options, but each use-case has its own Options as well, offering a further possibility for configuration.

#### ValidatorOption
You can customize the `Validator` during instantiation with Options. Using is just a matter of passing them as a constructor parameter, separated by a comma:
```go
service.NewValidator(service.WithJWKSURL("https://authservice.bitrise.io"), service.WithRealm("master"))
```

The available Options are the following:
- `WithBaseURL(url string)` overrides the authentication service's base URL.

- `WithSignatureAlgorithm(sa jose.SignatureAlgorithm)` overrides the encryption/decription algorithm of the *JWT*.
- `WithRealm(realm string)` overrides the realm.
- `WithKeyCacher(kc auth0.KeyCacher)` overrides the *JWK* cacher.
- `WithRealmURL(realmURL string)` overrides the realm URL.
- `WithJWKSURL(jwksURL string)` overrides the keystore URL.
- `WithValidator(validator JWTValidator)` overrides the Auth0 `Validator`.


#### HTTPMiddlewareOption

#### EchoMiddlewareOption


### Usage

#### Handler Function
```go
package main

import (
	"log"
	"net/http"

	"github.com/bitrise-io/bitrise-oauth/service"
)

func main() {
	handler := func(w http.ResponseWriter, r *http.Request) {}

	mux := http.NewServeMux()

	validator := service.NewValidator()

	mux.HandleFunc("/test_func", validator.HandlerFunc(handler))

	log.Fatal(http.ListenAndServe(":8080", mux))
}
```

#### Handler Function with gorilla/mux
```go
package main

import (
	"log"
	"net/http"

	"github.com/bitrise-io/bitrise-oauth/service"
	"github.com/gorilla/mux"
)

func main() {
	handler := func(w http.ResponseWriter, r *http.Request) {}

	router := mux.NewRouter()

	validator := service.NewValidator()

	router.HandleFunc("/test_func", validator.HandlerFunc(handler)).Methods(http.MethodGet)

	http.Handle("/", router)

	log.Fatal(http.ListenAndServe(":8080", router))
}

```

#### Middleware
```go
package main

import (
	"log"
	"net/http"

	"github.com/bitrise-io/bitrise-oauth/service"
)

func main() {
	handler := func(w http.ResponseWriter, r *http.Request) {}

	mux := http.NewServeMux()

	validator := service.NewValidator()

	mux.Handle("/test", validator.Middleware(http.HandlerFunc(handler)))

	log.Fatal(http.ListenAndServe(":8080", mux))
}
```

#### Middleware with gorilla/mux
```go
package main

import (
	"log"
	"net/http"

	"github.com/bitrise-io/bitrise-oauth/service"
	"github.com/gorilla/mux"
)

func main() {
	handler := func(w http.ResponseWriter, r *http.Request) {}

	router := mux.NewRouter()

	validator := service.NewValidator()

	router.Handle("/test", validator.Middleware(http.HandlerFunc(handler))).Methods(http.MethodGet)

	http.Handle("/", router)

	log.Fatal(http.ListenAndServe(":8080", router))
}
```

#### Echo Middleware Function
```go
package main

import (
	"net/http"

	"github.com/bitrise-io/bitrise-oauth/service"
	"github.com/labstack/echo"
)

func main() {
	handler := func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	}

	e := echo.New()

	validator := service.NewValidator()

	e.Use(validator.MiddlewareFunc())

	e.GET("/test", handler)

	e.Logger.Fatal(e.Start(":8080"))
}
```

#### Echo Handler Function
```go
package main

import (
	"net/http"

	"github.com/bitrise-io/bitrise-oauth/service"
	"github.com/labstack/echo"
)

func main() {
	validator := service.NewValidator()

	handler := func(c echo.Context) error {
		if err := validator.ValidateRequest(c.Request()); err != nil {
			return err
		}
		return c.String(http.StatusOK, "Hello, World!")
	}

	e := echo.New()

	e.GET("/test", handler)

	e.Logger.Fatal(e.Start(":8080"))
}
```








# Legacy
# Godoc

## Launch & view
1. Clone this project under `$GOPATH/src/github.com/bitrise-io/bitrise-oauth`
1. Run `godoc -http=:6060 &`
1. Open the documentation: http://localhost:6060/pkg/github.com/bitrise-io/bitrise-oauth

## Install godoc

```bash
go get golang.org/x/tools/cmd/godoc
```
