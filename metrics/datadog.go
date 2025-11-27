package metrics

import (
	"fmt"

	"github.com/DataDog/datadog-go/v5/statsd"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const unknownIssuerId = "unknown"

type DatadogMetrics struct {
	rawClient statsd.ClientInterface
	logger    Logger
}

// nolint: govet
type DatadogConfig struct {
	StatsdHost     string `env:"STATSD_HOST,default=datadog"`
	StatsdPort     string `env:"STATSD_PORT,default=8125"`
	MetricsEnabled bool   `env:"METRICS_ENABLED,default=true"`
}

type Logger interface {
	Errorw(msg string, keysAndValues ...interface{})
}

func NewDatadogMetrics(cfg DatadogConfig, logger Logger) (*DatadogMetrics, error) {
	if !cfg.MetricsEnabled {
		return &DatadogMetrics{
			rawClient: &statsd.NoOpClient{},
			logger:    logger,
		}, nil
	}

	c, err := statsd.New(fmt.Sprintf("%s:%s", cfg.StatsdHost, cfg.StatsdPort))
	if err != nil {
		return &DatadogMetrics{}, errors.Wrap(err, "create statsd client failed")
	}

	return &DatadogMetrics{
		rawClient: c,
		logger:    logger,
	}, nil
}

func (dm *DatadogMetrics) IncrRaw(name string, tags []string, rate float64) {
	if name == "" {
		dm.logger.Errorw("metric name is empty")

		return
	}

	if err := dm.rawClient.Incr(name, tags, rate); err != nil {
		dm.logger.Errorw("failed to increment raw metric", zap.String("metric", name), zap.Error(err))
	}
}

func (dm *DatadogMetrics) IncrAuthValidationSucceededMetric(issuer string) {
	if issuer == "" {
		issuer = unknownIssuerId
	}

	dm.IncrRaw("bitrise.jwt_auth.validation_succeeded", []string{"iss:" + issuer}, 1)
}

func (dm *DatadogMetrics) IncrAuthValidationFailedMetric(issuer string) {
	if issuer == "" {
		issuer = unknownIssuerId
	}

	dm.IncrRaw("bitrise.jwt_auth.validation_failed", []string{"iss:" + issuer}, 1)
}

func (dm *DatadogMetrics) Close() error {
	return dm.rawClient.Close()
}
