package config

// ...
const (
	BaseURL  = `http://104.154.234.133/auth`
	Realm    = `master`
	RealmURL = BaseURL + `/auth/realms/` + Realm
	JWKSURL  = RealmURL + `/protocol/openid-connect/certs`
	TokenURL = RealmURL + `/protocol/openid-connect/token`
)
