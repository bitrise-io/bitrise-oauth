package mocks

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
	"github.com/stretchr/testify/mock"
)

// token is signed with the private key below
var RawMockToken = `eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsImtpZCI6ImMwMGM5Yjk2MWU0OGM4YTkzMDYwOGY2NmQ2ODE3OTFiIn0.eyJTY29wZXMiOlsiZm9vIiwiYmFyIl0sImlzcyI6Imh0dHBzOi8vdG9rZW4taXNzdWVyLmJpdHJpc2UuaW8vYXV0aC9yZWFsbXMvYml0cmlzZS1zZXJ2aWNlcyIsInN1YiI6InN1YmplY3QiLCJhdWQiOiJ0ZXN0X2F1ZGllbmNlIn0.Ccb05Zu1HWvJCcs75kUghj_YyKVkbkm694wgrXA2pgkpMuEk6rPItZZd1UQcaBHzD2CAEPzVTQkx4wpvWG630tI2I1AkXoQmeEZJ52rygZIxbIOmxVzefxK_kDi-yl5SWnB5PBPMU4_0PKnObZYtDwt1MTjDZASON5xRoSKQRxY1MRdACNMB_-ayMfRzwoL76M7BzDvnLgQAbNrlsJkwWSuvqYFqL8995BqqChkxHndSShcjYZX8R8GVh0F1crmbb1J9-Twv5UmGPt4e9nYXaRxTTbMJwXaLkvy_Q-w9IKTjse99oMmNYRQ7CU3mnGlz5sdFJ6iHyDkawYnFCsxI0nBVE7NotcI18K-VK66s6coepwn4qaIR0pcicvzXUnSdmYwubRjoThFXz2iRXaeSr9sr-qInnZCVtatNVMUH4coH2XkF70QaArSpO4I3lrwHqiaFzAS8nI3TwNP66mxY7oCmPJHM-6ZBogLMBFTJDwjCsEGds1qBP9oTzccnbvOP`
var MockToken, _ = jwt.ParseSigned(RawMockToken, []jose.SignatureAlgorithm{jose.RS256})

// mock RS256 key pair in JWK (generated with https://www.scottbrady.io/tools/jwt)
// private-public key in JWK:
// {"alg":"RS256","d":"BKJqre6mb_BnULcjPL3GLuKkZiZNeAjVXohkML7-3Z4wzkzJiK5t_msjuqw6LWBUYFO0SW1FihV0t-ZZT5DTvrbY9awu77B_D2O1StVRHwslwCrLvAmwPoKZvn-lqhgGTJMZcFHef_aTH2o0LLC0qGrL5B2Fr8FkO8sXRvPwjnLrC7ehgTSOQ3aHH0NrUSsFdYpXUgXIh4kCLv3YDcgl8Rc1dkG9h3EIXCY1pNIdoL-ARMcE-mOll75uuQUvMHgh2D-KKgoT1eLGS1mglZE86QN9lv2L1qo4dy5PiXDT8S-uc3qiQVHC0vblgq_e9bwVPU63FbP0nrHZC-QDjmkwUjK52X8XRixClvHufb4_S8_eYSF803uDR55mDlZb-dCLGpIp0nkbxPUZyK1MjjuwI8yC4sSwuapgq1f85X194K3LeB5WMHULBi90WPQh2g4kSzjqCxy_HSNSby5kbYa6QshXX83iUeX23EMZqE9SWZs1RmhqMVhygAajXPqDt4qB","dp":"V_DWpg14_9GkWDWk8alWUoEuuX2ytjsCRnI0YhYLf-qSlBj6xkiCO0rLarWj4Fgo1Yngj5PrkoDyFWBe-LsDisrJn6njKLyC_TceDTmcwzoSBHcUN5JWj50CCdgL3d_UVoDG1aoSm1CJp8TF0meWb2hd7eNycsPPKXSsHAcobG0IQ193OR6EAchGLSiehLMZVv2iFTk-AGYJsUJS3yfPOQuEuazmpUGDf_gDlND1yFWV22BIIOaZNJ5cPoQzL_Yh","dq":"wZgR1Mhkxm--sRWx3ljlPstE4vJac3Hy8Cca0rTqGJ8cpU04PnxGKeOGZHsiXa90e3It040caRiQoqTaTS2RYJh5V7I3CYe-5TxScCRPrKGc0oEuoIhzTsv32Up77OjCK3KByGtBs5UWf8y06k8xf3M_jCIkUzehtHGgxFKbkYqDsfUumrvdY_QFw3WWvdE2DlDoRm8e7DGA9EyHkkjA5DzbylccpvZ0yTyWvGgOen_DDvFfxL4XKcZrB-58s_iB","e":"AQAB","key_ops":["sign"],"kty":"RSA","n":"wDO6jwSi57b6trpfeq0VnL80qBt6Z6XU4uw7UGwnJCzxTAiPHyN3P2drk_Bseb5uWN6v3zRcZbwrzCvGJvu3gwg_6YFVeYeDPeMIEPJYGu97sn2dZXCYT9c3n1l9oRlLjGWTcnZPwfdY0Ivmtfy8Q3cpAxq9BIIVfrIhTEvl-yLJ36DxMOSoUhmcF-bSc6Q_btDgRDf6tKLDnVi-IIN_vCWC7kDYn3-qAfSLpg0dGO2n_4-IWTN49J-lMQfE-dsKbyoVPj7w7ZVWDJEMTsJz_X89xsrKLOZ88BBp3HqmYsvFJRjgdJbYkYdPOHBe1Xv0-59-zzZv_RXL83jUz5EoYr5dQSBdaTRek-mxmIPo0q6bffRNimUumWsHN0d50aoiyvqK8m03tK7qlsWZ_oUAYMrUutjh5aTF6yekwDYbxYfRADmjlikaGu2uxp9G7w_Hs3bV0dyjCfofbWaVS73CI8aLMyCI8en6aa0eT7w86rtXoEOBmK_Vetz4gQig93XR","p":"-uXZt6xzSv-W5-T0PB3SjQ-3XiRmcHeQTcs6nLI9hGLba4zuh0OHufq_aYRCFi2MrGnDJSCtZAFYCKzr9p29-lefYeBnRZJEsQrhR0mb9ENTmcuEXs0ax5fLkEQtwiUL8H_R45d3vZ0owkDP0C7-mIKVbZlc760bQyp9NHiiCUQiSPLju6MF4GAFO0lQ4-aGLjz46Rcu28rAriqNFVlqm0XzjcPNEko1WqH4DVW4CxfF-GDzlkqf6ibnZD_mMPBh","q":"xBxQVG9GrsHxgUSL3v52jK4fA_xRUJzpP3I5x2YDiSMwvFPQkA0v62YXRFnNgyguAdzFrl8X_OrcXkriqcTstuHvSVQq6InMlEGvM1c1-g3Gp_UdajunlSCZnU5-x2LrqXjBdbdLjsM6rWw_6Le2ndq0oA7hSJIUPoBqk6fetu5TPZsFOQ2s_3eqXoFy9JqPBM6881XXBNhS9gS7ZbYQk2FXbhQY7-ZHWdAStmVHNDF_-qLVAQlCSV5STGQtVztx","qi":"iPEDQcxmYysdVJ19Ft9MWvrY95sVo0GKToOWaBT0VJkMHhOUwCrMG-huCs-wofcFLdtxLz9QxiAwa1OKv2OPW4tgDIQu__aEDtrWpJ-bmDdBgFEUpuUh31jQ-7OJB6eM1duL9yg0mZvwYU4tsDPTIoYIS0GfvKiaZNCzp_ORd3CMVymA4cmRu-Zp2f7GhYi17M8nVFs7i4sKCV-q5ApeQh4IxRl_xVgUy7-tSYpBi6wny4JF1JtupxwKLNzPmm32","use":"sig","kid":"c00c9b961e48c8a930608f66d681791b"}
const mockRawPublickKeyJWK = `{"alg":"RS256","e":"AQAB","key_ops":["verify"],"kty":"RSA","n":"wDO6jwSi57b6trpfeq0VnL80qBt6Z6XU4uw7UGwnJCzxTAiPHyN3P2drk_Bseb5uWN6v3zRcZbwrzCvGJvu3gwg_6YFVeYeDPeMIEPJYGu97sn2dZXCYT9c3n1l9oRlLjGWTcnZPwfdY0Ivmtfy8Q3cpAxq9BIIVfrIhTEvl-yLJ36DxMOSoUhmcF-bSc6Q_btDgRDf6tKLDnVi-IIN_vCWC7kDYn3-qAfSLpg0dGO2n_4-IWTN49J-lMQfE-dsKbyoVPj7w7ZVWDJEMTsJz_X89xsrKLOZ88BBp3HqmYsvFJRjgdJbYkYdPOHBe1Xv0-59-zzZv_RXL83jUz5EoYr5dQSBdaTRek-mxmIPo0q6bffRNimUumWsHN0d50aoiyvqK8m03tK7qlsWZ_oUAYMrUutjh5aTF6yekwDYbxYfRADmjlikaGu2uxp9G7w_Hs3bV0dyjCfofbWaVS73CI8aLMyCI8en6aa0eT7w86rtXoEOBmK_Vetz4gQig93XR","use":"sig","kid":"c00c9b961e48c8a930608f66d681791b"}`

var MockPublicKey = func() jose.JSONWebKey {
	var jwk = jose.JSONWebKey{}
	err := json.NewDecoder(strings.NewReader(mockRawPublickKeyJWK)).Decode(&jwk)
	if err != nil {
		panic(err)
	}

	return jwk
}

// JWTValidator ...
type JWTValidator struct {
	mock.Mock
}

// ValidateRequest ...
func (m *JWTValidator) ValidateRequest(r *http.Request) (*jwt.JSONWebToken, error) {
	args := m.Called(r)
	return args.Get(0).(*jwt.JSONWebToken), args.Error(1)
}

// GivenSuccessfulJWTValidation ...
func (m *JWTValidator) GivenSuccessfulJWTValidation() *JWTValidator {
	m.On("ValidateRequest", mock.Anything).Return(MockToken, nil)
	return m
}

// GivenUnsuccessfulJWTValidation ...
func (m *JWTValidator) GivenUnsuccessfulJWTValidation(err error) *JWTValidator {
	m.On("ValidateRequest", mock.Anything).Return(&jwt.JSONWebToken{}, err)
	return m
}
