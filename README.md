
# Bitrise OAuth library for Go

## Client
lorem ipsum intro




## Server
lorem ipsum intro

## Usage

### Validator as Handler Function
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

### Validator as Handler Function with gorilla/mux
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

### Validator as Middleware
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

### Validator as Middleware with gorilla/mux
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

### Validator as Echo Middleware Function
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

### Validator as Echo Handler Function
```go
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
