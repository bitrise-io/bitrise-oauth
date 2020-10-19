package config

// ...
const (
	BaseURL  = `https://auth.services.bitrise.io`
	Realm    = `master`
	RealmURL = BaseURL + `/auth/realms/` + Realm
	JWKSURL  = RealmURL + `/protocol/openid-connect/certs`
	TokenURL = RealmURL + `/protocol/openid-connect/token`
)
