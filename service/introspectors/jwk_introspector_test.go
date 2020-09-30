package introspectors_test

import (
	"log"
	"net/http"

	"github.com/bitrise-io/bitriseoauth/service/introspectors"
	"github.com/labstack/echo"
)

func ExampleJWK_Middleware() {
	handler := func(w http.ResponseWriter, r *http.Request) {}

	mux := http.NewServeMux()

	introspector := introspectors.NewJWK(nil, nil, nil)

	mux.Handle("/test", introspector.Middleware(http.HandlerFunc(handler)))

	log.Fatal(http.ListenAndServe(":8080", mux))
}

func ExampleJWK_HandlerFunc() {
	handler := func(w http.ResponseWriter, r *http.Request) {}

	mux := http.NewServeMux()

	introspector := introspectors.NewJWK(nil, nil, nil)

	mux.HandleFunc("/test_func", introspector.HandlerFunc(handler))

	log.Fatal(http.ListenAndServe(":8080", mux))
}

func ExampleJWK_ValidateRequest() {
	introspector := introspectors.NewJWK(nil, nil, nil)

	handler := func(c echo.Context) error {
		if err := introspector.ValidateRequest(c.Request()); err != nil {
			return err
		}
		return c.String(http.StatusOK, "Hello, World!")
	}

	e := echo.New()

	e.GET("/test", handler)

	e.Logger.Fatal(e.Start(":1323"))
}

func ExampleJWK_MiddlewareFunc_echo() {
	handler := func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	}

	e := echo.New()

	introspector := introspectors.NewJWK(nil, nil, nil)

	e.Use(introspector.MiddlewareFunc())

	e.GET("/test", handler)

	e.Logger.Fatal(e.Start(":1323"))
}
