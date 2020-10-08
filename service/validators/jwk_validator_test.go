package validators_test

import (
	"log"
	"net/http"

	"github.com/bitrise-io/bitrise-oauth/service/validators"
	"github.com/labstack/echo"
)

func ExampleJWK_Middleware() {
	handler := func(w http.ResponseWriter, r *http.Request) {}

	mux := http.NewServeMux()

	validator := validators.NewJWK(nil, nil, nil)

	mux.Handle("/test", validator.Middleware(http.HandlerFunc(handler)))

	log.Fatal(http.ListenAndServe(":8080", mux))
}

func ExampleJWK_HandlerFunc() {
	handler := func(w http.ResponseWriter, r *http.Request) {}

	mux := http.NewServeMux()

	validator := validators.NewJWK(nil, nil, nil)

	mux.HandleFunc("/test_func", validator.HandlerFunc(handler))

	log.Fatal(http.ListenAndServe(":8080", mux))
}

func ExampleJWK_ValidateRequest() {
	validator := validators.NewJWK(nil, nil, nil)

	handler := func(c echo.Context) error {
		if err := validator.ValidateRequest(c.Request()); err != nil {
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

	validator := validators.NewJWK(nil, nil, nil)

	e.Use(validator.MiddlewareFunc())

	e.GET("/test", handler)

	e.Logger.Fatal(e.Start(":1323"))
}
