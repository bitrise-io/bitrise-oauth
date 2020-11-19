package service

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

func (audienceConfig AudienceConfig) all() []string {
	return audienceConfig.audience
}
