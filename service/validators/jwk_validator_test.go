package validators_test

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	auth0 "github.com/auth0-community/go-auth0"
	"github.com/bitrise-io/bitrise-oauth/mocks"
	"github.com/bitrise-io/bitrise-oauth/service"
	"github.com/bitrise-io/bitrise-oauth/service/validators"
	"github.com/gorilla/mux"
	"github.com/labstack/echo"
)

const (
	Authorization = "Authorization"
	// You may find further information about generating signed JWT tokens here: https://bitrise.atlassian.net/wiki/spaces/~940361272/pages/828867013/Generate+JWT+for+testing
	JWT_1      = "BEARER eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IldNVlkwcWVyZmJGUC0zcllLdW55NUFQaXJmcnk0OG5QZWVYcnlQNzk5RmsifQ.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.TG_KUN36-raydoUZH0IeF_-upVbRCLemD7Bt3BCWCrr51cJjsSKAkIKmnyFH4Ow_5pLzlaU4NRoQMeoAIF0VvH28P6hnQKvMmpEM-kQ0UMHnZvMfzuh7SvsvEAIaDhtEEOFfYNV5w0jWoQsAbrMw9vkKrPFqXatxBF1t_yvbW6x2SKEql_UmmN88oxfa_-DS2OrKWsyk2hakR6rnM-m8zTrqdsRndvAP25DeNiySHR_fyB53Dn7un-TO7KJENi5X_obGuXKjQY0C5JFkibR1RY4o9Rp04rdFQrv_PPBC2Ki0pIqDpKVaRbceNkC1BiMzz2zjNR2B6EWbdzPB24bgAA"
	JWT_2      = "BEARER eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IldNVlkwcWVyZmJGUC0zcllLdW55NUFQaXJmcnk0OG5QZWVYcnlQNzk5RmsxIn0.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.XLWuYJ3xb71XGh3XXH0xpX18_Q8RWQDIjUs6EKYD1mX2KXkJIWKj_1m4gNziEdTm03tXFKuDCXhdnFB7L7NJeOtT9dVNtIfqkBR0cYd2KU6HitPck9qd5wz_JcaaPQttHfrVBxJVIaK7ifZFCjjlGSukkYQ7aQalEv2ZjTycXP7FVs7bDq39f1OWdw2rM6XurrjWm65uEwC9m2z08DdgPnmyzCFh0NE5WyMHkezcIl2DDHxJjmb0AZkdIYW1q-AbYs0CIlAemOnxW_or7uzgtATZ-GWE_WEJp_bOeTkZK3BLnShXhlRdKNaHJXCuBzfBwdUY24-x6mEPRKNBYPGW3w"
	RequestUrl = "https://bitrise.io/protected_route"
)

func ExampleJWK_Middleware() {
	handler := func(w http.ResponseWriter, r *http.Request) {}

	mux := http.NewServeMux()

	validator := validators.NewJWK(nil, nil, nil)

	mux.Handle("/test", validator.Middleware(http.HandlerFunc(handler)))

	log.Fatal(http.ListenAndServe(":8080", mux))
}

func ExampleJWK_Middleware_gorilla_mux() {
	handler := func(w http.ResponseWriter, r *http.Request) {}

	router := mux.NewRouter()

	validator := validators.NewJWK(nil, nil, nil)

	router.Handle("/test", validator.Middleware(http.HandlerFunc(handler))).Methods(http.MethodGet)

	http.Handle("/", router)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func ExampleJWK_HandlerFunc() {
	handler := func(w http.ResponseWriter, r *http.Request) {}

	mux := http.NewServeMux()

	validator := validators.NewJWK(nil, nil, nil)

	mux.HandleFunc("/test_func", validator.HandlerFunc(handler))

	log.Fatal(http.ListenAndServe(":8080", mux))
}

func ExampleJWK_HandlerFunc_with_gorilla_mux() {
	handler := func(w http.ResponseWriter, r *http.Request) {}

	router := mux.NewRouter()

	validator := validators.NewJWK(nil, nil, nil)

	router.HandleFunc("/test_func", validator.HandlerFunc(handler)).Methods(http.MethodGet)

	http.Handle("/", router)

	log.Fatal(http.ListenAndServe(":8080", router))
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

func Test_Auth0_JWKS_Caching(t *testing.T) {
	testCases := []struct {
		name         string
		token1       string
		token2       string
		expiryInSecs time.Duration
		want         int
	}{
		{
			"1. Given two requests with the same token that will NOT expire when the requests are validated then expect /certs endpoint to be called ONCE",
			JWT_1,
			JWT_1,
			60,
			1,
		},
		{
			"2. Given two requests with different tokens that will NOT expire when the requests are validated then expect /certs endpoint to be called TWICE",
			JWT_1,
			JWT_2,
			60,
			2,
		},
		{
			"3. Given two requests with the same token that will expire when the requests are validated then expect /certs endpoint to be called TWICE",
			JWT_1,
			JWT_1,
			1,
			2,
		},
		{
			"4. Given two requests with different token that will expire when the requests are validated then expect /certs endpoint to be called TWICE",
			JWT_1,
			JWT_2,
			1,
			2,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// Given
			mockAuthService := mocks.AuthService{}
			mockAuthService.On("Certs").Return().Times(testCase.want)

			testAuthServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Path {
				case "/certs":
					mockAuthService.Certs()
					addContentTypeAndTokenToResponse(w)
				default:
					fmt.Println("Request handled by the default endpoint")
					w.WriteHeader(http.StatusOK)
				}
			}))
			defer testAuthServer.Close()

			validator := validators.NewJWK(validators.WithJWKSURL(testAuthServer.URL+"/certs"), validators.WithKeyCacher(auth0.NewMemoryKeyCacher(testCase.expiryInSecs*time.Millisecond, 5)))

			request1 := createRequestWithToken(testCase.token1)
			request2 := createRequestWithToken(testCase.token2)

			// When
			validateRequest(validator, request1)
			time.Sleep(20 * time.Millisecond)
			validateRequest(validator, request2)

			// Then
			mockAuthService.AssertExpectations(t)
		})
	}
}

func addContentTypeAndTokenToResponse(w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
	_, err := w.Write([]byte("{\"keys\": [{\"kid\": \"WMVY0qerfbFP-3rYKuny5APirfry48nPeeXryP799Fk\",\"kty\": \"RSA\",\"alg\": \"RS256\",\"use\": \"sig\",\"n\": \"kGb4ABWuOgQH7yCydKsqLjZ7-FrWOQ5tezQQofHs5jJYQPXnMalUgvY6v9c0GEBvzebbmkigcGw9e8NesOLnVaP4mkE6TYLDyuL1vDRP9bIQuVOQxDwqhDPmaFKFawxe0YLoFN_N6NOZBZJ69z2Mbhxsd9By4tr_-bR-seg-korL5NUf6KpLYdBeCDy1xK_DSia49vV-SYAG5cuxgejRc32fmmZnVFx7rs8nIIAUoaGAHhWGM7ZFRaxC96dFsVRRR85-TDeukkPi0_-Q2RtNwpoz-hP2g-p6Vl620z3KJYW6pO6ssT3Q8SY0LOoK6gPca7NW7qGvOGyWhM6yk4aqPQ\",\"e\": \"AQAB\",\"x5c\": [\"MIICmzCCAYMCBgF01MwbDDANBgkqhkiG9w0BAQsFADARMQ8wDQYDVQQDDAZtYXN0ZXIwHhcNMjAwOTI4MTI1ODAwWhcNMzAwOTI4MTI1OTQwWjARMQ8wDQYDVQQDDAZtYXN0ZXIwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCQZvgAFa46BAfvILJ0qyouNnv4WtY5Dm17NBCh8ezmMlhA9ecxqVSC9jq/1zQYQG/N5tuaSKBwbD17w16w4udVo/iaQTpNgsPK4vW8NE/1shC5U5DEPCqEM+ZoUoVrDF7RgugU383o05kFknr3PYxuHGx30HLi2v/5tH6x6D6Sisvk1R/oqkth0F4IPLXEr8NKJrj29X5JgAbly7GB6NFzfZ+aZmdUXHuuzycggBShoYAeFYYztkVFrEL3p0WxVFFHzn5MN66SQ+LT/5DZG03CmjP6E/aD6npWXrbTPcolhbqk7qyxPdDxJjQs6grqA9xrs1buoa84bJaEzrKThqo9AgMBAAEwDQYJKoZIhvcNAQELBQADggEBACSleX2IEZqb4h05T9+W3BT2e8cRiX06T8pHk4o70GAooROvMeHu0+l+HlT6lYggerzxsxqYGA7KFOG6JmgFG5XLPlNJoHlX0NOCGfCbrh50Q1HV5TZsqubIUOHglos9/SotipiSSVncd02Uot27sqU3A1HR9qO3IxmTe+W5XIvmzn3Kofpyj5r9qzbSJMfW2YKCJ8n+lG0g184SQ1JtQ2zFwBdziHtO8eBTscadjnZy4WHTs7F9hh2xNGUVBsfDQ5JZooJaZEeGe9r9Fv46R5Py/SMZvcpp6PvNptN2ifXoPzcP6jAZVphRlv7DZXIjb+9UgN3fHcDZNTK3NvJWoDA=\"],\"x5t\": \"fp4lFC6SkVBMvUGHq2MgrvB20L0\",\"x5t#S256\": \"BN6wsc3PZ_XfJXKPqH-SGdhe7QmEz-fECiheWPvkEBA\"}]}"))
	if err != nil {
		panic("Can't write body")
	}
}

func createRequestWithToken(jwt string) *http.Request {
	request, err := http.NewRequest(http.MethodGet, RequestUrl, nil)
	if err != nil {
		panic("Can't create request")
	}

	request.Header.Add(Authorization, jwt)
	return request
}

func validateRequest(validator service.Validator, request *http.Request) {
	err := validator.ValidateRequest(request)
	if err != nil {
		panic("Can't validate request")
	}
}
