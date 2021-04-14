package service

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/auth0-community/go-auth0"
	"github.com/bitrise-io/bitrise-oauth/config"
	"github.com/bitrise-io/bitrise-oauth/mocks"
	"github.com/c2fo/testify/assert"
	"github.com/c2fo/testify/mock"
	"github.com/labstack/echo"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

func Test_GivenSuccessfulJWTValidationWithMiddleware_WhenRequestIsHandled_ThenExpectTheNextMiddlewareToBeCalled(t *testing.T) {
	// Given
	mockHandler := givenMockHandler()
	mockErrorWriter := givenMockErrorWriter()
	mockSecretProvider := givenMockSecretProvider()

	validator := createValidator(givenSuccessfulJWTValidation(), mockSecretProvider)
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

	validator := createValidator(givenUnsuccessfulJWTValidation(), nil)
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
	mockSecretProvider := givenMockSecretProvider()

	validator := createValidator(givenSuccessfulJWTValidation(), mockSecretProvider)
	validatorMiddlewareFunction := validator.EchoMiddlewareFunc(WithContextErrorWriter(mockErrorWriter.EchoHandlerFunc))(mockMiddlewareHandlerFunction.HandlerFunction)

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

	validator := createValidator(givenUnsuccessfulJWTValidation(), nil)
	validatorMiddlewareFunction := validator.EchoMiddlewareFunc(WithContextErrorWriter(mockErrorWriter.EchoHandlerFunc))(mockMiddlewareHandlerFunction.HandlerFunction)

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
	mockSecretProvider := givenMockSecretProvider()

	validator := createValidator(givenSuccessfulJWTValidation(), mockSecretProvider)
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

	validator := createValidator(givenUnsuccessfulJWTValidation(), nil)
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

func givenMockSecretProvider() *mocks.MockSecretProvider {
	mockSecretProvider := new(mocks.MockSecretProvider)
	mockSecretProvider.On("GetSecret", mock.Anything).Return([]byte("secret"), nil)
	return mockSecretProvider
}

func givenMockEchoErrorWriter(err error) *mocks.ErrorWriter {
	mockErrorWriter := new(mocks.ErrorWriter)
	mockErrorWriter.On("EchoHandlerFunc", mock.Anything, mock.Anything).Return(err)
	return mockErrorWriter
}

// type mockValidator struct {
// 	mock.Mock
// }

// func (m *mockValidator) ValidateRequest(r *http.Request) error {
// 	return nil
// }

// func (m *mockValidator) HandlerFunc(http.HandlerFunc, ...HTTPMiddlewareOption) http.HandlerFunc {
// 	return nil
// }
// func (m *mockValidator) Middleware(http.Handler, ...HTTPMiddlewareOption) http.Handler {
// 	return nil
// }
// func (m *mockValidator) EchoMiddlewareFunc(...EchoMiddlewareOption) echo.MiddlewareFunc {
// 	return nil
// }
// func (m *mockValidator) ValidateRequestAndReturnToken(r *http.Request) (TokenWithClaims, error) {
// 	return nil, nil
// }
// func (m *mockValidator) ValidateAudiences(tokenWithClaims tokenWithClaims, audiences []string) error {
// 	return nil
// }

// func (m *mockValidator) GivenSuccessfulValidateRequest() {
// 	m.On("ValidateRequest").Return(nil)
// }

func createValidator(mockJWTValidator jwtValidator, mockSecretProvider auth0.SecretProvider) Validator {
	validator := NewValidator(
		config.NewAudienceConfig("test_audience"),
		withValidator(mockJWTValidator),
		withSecretProvider(mockSecretProvider),
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

func Test_AudienceClaimValidation(t *testing.T) {
	testCases := []struct {
		name           string
		tokenAudiences []string
		inputAudiences []string
		expectedError  error
	}{
		{
			name:           "Given a request and a validator with the SAME audience when the request is vaidated then expect no error",
			tokenAudiences: []string{"aud1"},
			inputAudiences: []string{"aud1"},
			expectedError:  nil,
		},
		{
			name:           "Given a request and a validator with a common subset of audiences",
			tokenAudiences: []string{"aud1", "aud2", "aud3"},
			inputAudiences: []string{"aud2", "aud3", "aud4"},
			expectedError:  nil,
		},
		{
			name:           "Given a request and a validator without a common subset of audiences",
			tokenAudiences: []string{"aud1", "aud2", "aud3"},
			inputAudiences: []string{"aud4", "aud5"},
			expectedError:  jwt.ErrInvalidAudience,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// Given
			testToken := newTestTokenConfigWithAudiences(testCase.tokenAudiences)
			request := testToken.newRequest()

			validator := NewValidator(
				config.NewAudienceConfig(testCase.inputAudiences[0], testCase.inputAudiences[1:]...),
				withIssuer(defaultIssuer),
				withSecretProvider(defaultSecretProvider),
			)

			// When
			err := validator.ValidateRequest(request)

			// Then
			assert.Equal(t, testCase.expectedError, err)
		})
	}
}

func newTestTokenConfigWithAudiences(audiences []string) testTokenConfig {
	return testTokenConfig{
		audiences,
		defaultIssuer,
		time.Now().Add(24 * time.Hour),
		jose.RS256,
		defaultSecret,
		defaultKid,
	}
}
