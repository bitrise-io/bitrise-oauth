package config

// AudienceConfig ...
type AudienceConfig struct {
	audience []string
}

// NewAudienceConfig ...
func NewAudienceConfig(audience string, audiences ...string) AudienceConfig {
	return AudienceConfig{
		audience: append(audiences, audience),
	}
}

// NewAudienceConfigFromAudiences creates a new AudienceConfig from a slice of audiences
func NewAudienceConfigFromAudiences(audiences []string) AudienceConfig {
	return AudienceConfig{
		audience: audiences,
	}
}

// All ...
func (audienceConfig AudienceConfig) All() []string {
	return audienceConfig.audience
}
