package newrelic

import (
	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
)

// Metric describes a single fetchable metric.
type Metric struct {
	Namespace plugin.Namespace
	Type      string
	Path      string
	Unit      string
}

// Service is the interface every New Relic service component must implement.
type Service interface {
	GetMetricTypes(plugin.Config) ([]plugin.Metric, error)
	CollectMetrics([]plugin.Metric) ([]plugin.Metric, error)
}

// Collector takes care of fetching data from New Relic.
type Collector struct{}

// GetConfigPolicy defines the configuration variables this plugin supports.
func (n *Collector) GetConfigPolicy() (plugin.ConfigPolicy, error) {
	p := plugin.NewConfigPolicy()
	p.AddNewStringRule(
		[]string{"inteleon", "newrelic"},
		"api_key",
		true,
	)

	return *p, nil
}

// GetMetricTypes returns all the metric types this collector supports.
func (n *Collector) GetMetricTypes(cfg plugin.Config) ([]plugin.Metric, error) {
	ret := []plugin.Metric{}

	apm := &APM{}

	apmMet, err := apm.GetMetricTypes(cfg)
	if err != nil {
		return ret, err
	}

	for _, m := range apmMet {
		ret = append(ret, m)
	}

	return ret, nil
}

// CollectMetrics fetches all the requested metrics and returns them.
func (n *Collector) CollectMetrics(metrics []plugin.Metric) ([]plugin.Metric, error) {
	ret := []plugin.Metric{}

	cfg := metrics[0].Config
	apm := NewAPM(cfg["api_key"].(string))

	apmMet, err := apm.CollectMetrics(metrics)
	if err != nil {
		return ret, err
	}

	for _, m := range apmMet {
		ret = append(ret, m)
	}

	return ret, nil
}
