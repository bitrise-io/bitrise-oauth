package config

// ...
const (
	BaseURL  = `http://35.184.90.188`
	Realm    = `master`
	RealmURL = BaseURL + `/auth/realms/` + Realm
	JWKSURL  = RealmURL + `/protocol/openid-connect/certs`
	TokenURL = RealmURL + `/protocol/openid-connect/token`
)
