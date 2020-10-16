package config

// ...
const (
	BaseURL  = `http://35.232.43.235`
	Realm    = `master`
	RealmURL = BaseURL + `/auth/realms/` + Realm
	JWKSURL  = RealmURL + `/protocol/openid-connect/certs`
	TokenURL = RealmURL + `/protocol/openid-connect/token`
)
