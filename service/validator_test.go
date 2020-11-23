package service

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bitrise-io/bitrise-oauth/mocks"
	"github.com/c2fo/testify/assert"
	"github.com/c2fo/testify/mock"
	"github.com/labstack/echo"
	"gopkg.in/square/go-jose.v2"
)

func Test_GivenSuccessfulJWTValidationWithMiddleware_WhenRequestIsHandled_ThenExpectTheNextMiddlewareToBeCalled(t *testing.T) {
	// Given
	mockHandler := givenMockHandler()
	mockErrorWriter := givenMockErrorWriter()

	validator := createValidator(givenSuccessfulJWTValidation())
	testServer := startServerWithMiddleware(mockHandler, validator, WithHTTPErrorWriter(mockErrorWriter.ErrorHandler))

	// When
	sendGetRequest(testServer.URL)

	// Then
	testServer.Close()
	mockErrorWriter.AssertNotCalled(t, "ErrorHandler", mock.Anything, mock.Anything, mock.Anything)
	mockHandler.AssertCalled(t, "ServeHTTP", mock.Anything, mock.Anything)
}

func Test_GivenUnsuccessfulJWTValidationWithMiddleware_WhenRequestIsHandled_ThenExpectAnError(t *testing.T) {
	// Given
	mockHandler := givenMockHandler()
	mockErrorWriter := givenMockErrorWriter()

	validator := createValidator(givenUnsuccessfulJWTValidation())
	testServer := startServerWithMiddleware(mockHandler, validator, WithHTTPErrorWriter(mockErrorWriter.ErrorHandler))

	// When
	sendGetRequest(testServer.URL)

	// Then
	testServer.Close()
	mockErrorWriter.AssertCalled(t, "ErrorHandler", mock.Anything, mock.Anything, mock.Anything)
	mockHandler.AssertNotCalled(t, "ServeHTTP", mock.Anything, mock.Anything)
}

func Test_GivenSuccessfulJWTValidationWithMiddlewareHandlerFunction_WhenRequestIsHandled_ThenExpectTheNextMiddlewareToBeCalled(t *testing.T) {
	// Given
	mockMiddlewareHandlerFunction := givenMockMiddlewareHandlerFunctionWithSuccess()
	mockErrorWriter := givenMockErrorWriter()

	validator := createValidator(givenSuccessfulJWTValidation())
	validatorMiddlewareFunction := validator.MiddlewareFunc(WithContextErrorWriter(mockErrorWriter.EchoHandlerFunc))(mockMiddlewareHandlerFunction.HandlerFunction)

	context := createContext()

	// When
	err := validatorMiddlewareFunction(context)

	// Then
	assert.NoError(t, err)
	mockMiddlewareHandlerFunction.AssertCalled(t, "HandlerFunction", mock.Anything)
}

func Test_GivenUnsuccessfulJWTValidationWithMiddlewareHandlerFunction_WhenRequestIsHandled_ThenExpectAnError(t *testing.T) {
	// Given
	mockMiddlewareHandlerFunction := givenMockMiddlewareHandlerFunctionWithSuccess()
	mockErrorWriter := givenMockEchoErrorWriter(errors.New("error"))

	validator := createValidator(givenUnsuccessfulJWTValidation())
	validatorMiddlewareFunction := validator.MiddlewareFunc(WithContextErrorWriter(mockErrorWriter.EchoHandlerFunc))(mockMiddlewareHandlerFunction.HandlerFunction)

	context := createContext()

	// When
	err := validatorMiddlewareFunction(context)

	// Then
	assert.Error(t, err)
	mockMiddlewareHandlerFunction.AssertNotCalled(t, "HandlerFunction", mock.Anything)
}

func Test_GivenSuccessfulJWTValidationWithHandlerFunction_WhenRequestIsHandled_ThenExpectTheNextHandlerFunctionToBeCalled(t *testing.T) {
	//Given
	mockHandlerFunction := givenMockHandlerFunction()
	mockErrorWriter := givenMockErrorWriter()

	validator := createValidator(givenSuccessfulJWTValidation())
	testServer := startServerWithHandlerFunction(mockHandlerFunction.Handler, validator, WithHTTPErrorWriter(mockErrorWriter.ErrorHandler))

	// When
	sendGetRequest(testServer.URL)

	// Then
	testServer.Close()
	mockErrorWriter.AssertNotCalled(t, "ErrorHandler", mock.Anything)
	mockHandlerFunction.AssertCalled(t, "Handler", mock.Anything, mock.Anything)
}

func Test_GivenUnsuccessfulJWTValidationWithHandlerFunction_WhenRequestIsHandled_ThenExpectAnError(t *testing.T) {
	//Given
	mockHandlerFunction := givenMockHandlerFunction()
	mockErrorWriter := givenMockErrorWriter()

	validator := createValidator(givenUnsuccessfulJWTValidation())
	testServer := startServerWithHandlerFunction(mockHandlerFunction.Handler, validator, WithHTTPErrorWriter(mockErrorWriter.ErrorHandler))

	// When
	sendGetRequest(testServer.URL)

	// Then
	testServer.Close()
	mockErrorWriter.AssertCalled(t, "ErrorHandler", mock.Anything, mock.Anything, mock.Anything)
	mockHandlerFunction.AssertNotCalled(t, "Handler", mock.Anything, mock.Anything)
}

func givenSuccessfulJWTValidation() *mocks.JWTValidator {
	return new(mocks.JWTValidator).GivenSuccessfulJWTValidation()
}

func givenUnsuccessfulJWTValidation() *mocks.JWTValidator {
	return new(mocks.JWTValidator).GivenUnsuccessfulJWTValidation(errors.New("Can't validate request"))
}

func givenMockHandler() *mocks.Handler {
	mockHandler := new(mocks.Handler)
	mockHandler.On("ServeHTTP", mock.Anything, mock.Anything).Return()
	return mockHandler
}

func givenMockMiddlewareHandlerFunctionWithSuccess() *mocks.MiddlewareHandlerFunction {
	mockMiddlewareHandlerFunction := new(mocks.MiddlewareHandlerFunction)
	mockMiddlewareHandlerFunction.GivenSuccess()
	return mockMiddlewareHandlerFunction
}

func givenMockHandlerFunction() *mocks.HandlerFunction {
	mockHandlerFunction := new(mocks.HandlerFunction)
	mockHandlerFunction.On("Handler", mock.Anything, mock.Anything).Return()
	return mockHandlerFunction
}

func givenMockErrorWriter() *mocks.ErrorWriter {
	mockErrorWriter := new(mocks.ErrorWriter)
	mockErrorWriter.On("ErrorHandler", mock.Anything, mock.Anything, mock.Anything).Return()
	return mockErrorWriter
}

func givenMockEchoErrorWriter(err error) *mocks.ErrorWriter {
	mockErrorWriter := new(mocks.ErrorWriter)
	mockErrorWriter.On("EchoHandlerFunc", mock.Anything, mock.Anything).Return(err)
	return mockErrorWriter
}

func createValidator(mockJWTValidator jwtValidator) Validator {
	validator := NewValidator(
		NewAudienceConfig("test_audience"),
		withValidator(mockJWTValidator),
	)
	return validator
}

func startServerWithMiddleware(mockHandler *mocks.Handler, validator Validator, opts ...HTTPMiddlewareOption) *httptest.Server {
	testServer := httptest.NewServer(validator.Middleware(mockHandler, opts...))
	return testServer
}

func startServerWithHandlerFunction(mockHandlerFunction func(http.ResponseWriter, *http.Request), validator Validator, opt HTTPMiddlewareOption) *httptest.Server {
	testServer := httptest.NewServer(validator.HandlerFunc(mockHandlerFunction, opt))
	return testServer
}

func sendGetRequest(url string) {
	_, err := http.Get(url)
	if err != nil {
		panic(err)
	}
}

func createContext() echo.Context {
	echo := echo.New()

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()
	context := echo.NewContext(request, recorder)

	return context
}

func Test_qwe(t *testing.T) {
	token := getTestTokenWithKid(defaultAudience, defaultIssuer, time.Now().Add(24*time.Hour), jose.RS256, defaultSecret, defaultKid)
	req := createRequestWithToken(token)

	validator := NewValidator(
		NewAudienceConfig(defaultAudience[0]),
		withIssuer(defaultIssuer),
		withSecretProvider(defaultSecretProvider),
	)

	err := validator.ValidateRequest(req)

	fmt.Println(err)
}
