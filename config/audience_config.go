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

// All() ...
func (audienceConfig AudienceConfig) All() []string {
	return audienceConfig.audience
}
