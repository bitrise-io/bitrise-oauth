
# Bitrise OAuth library for Go

## Client
lorem ipsum intro




## Server
lorem ipsum intro

## Usage
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

### Validator as Handler Function

### Validator as Handler Function with gorilla/mux

### Validator as Echo Middleware Function







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
