package newrelic

import (
	"fmt"
	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
	"strings"
	"time"
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

	for _, comp := range []Service{NewAPM(""), NewCustom("")} {
		met, err := comp.GetMetricTypes(cfg)
		if err != nil {
			return ret, err
		}

		for _, m := range met {
			ret = append(ret, m)
		}
	}

	return ret, nil
}

// CollectMetrics fetches all the requested metrics and returns them.
func (n *Collector) CollectMetrics(metrics []plugin.Metric) ([]plugin.Metric, error) {
	ret := []plugin.Metric{}

	cfg := metrics[0].Config

	for _, comp := range []Service{NewAPM(cfg["api_key"].(string)), NewCustom(cfg["api_key"].(string))} {
		met, err := comp.CollectMetrics(metrics)
		if err != nil {
			return ret, err
		}

		for _, m := range met {
			ret = append(ret, m)
		}
	}

	return ret, nil
}

func populateMetric(metric plugin.Metric, mapData map[string]interface{}) (plugin.Metric, error) {
	// Create a new metric based on the "old" one.
	newMetric := metric

	mPath := []string{}
	if metric.Tags["Path"] == "" {
		mPath = []string{metric.Namespace.Element(len(metric.Namespace) - 2).Value}
	} else {
		mPath = strings.Split(metric.Tags["Path"], "/")
	}

	metricData, err := mapTraverse(mapData, mPath)
	if err != nil {
		return newMetric, err
	}

	newMetric.Data = metricData
	newMetric.Unit = metric.Tags["Unit"]
	newMetric.Tags = map[string]string{}
	newMetric.Timestamp = time.Now().UTC()

	return newMetric, nil
}

func mapTraverse(mapData map[string]interface{}, path []string) (interface{}, error) {
	pathElemNotFoundErrTemplate := "Path element not found: %s"

	if len(path) > 1 {
		pathElem, newPath := path[0], path[1:]

		newMapData, ok := mapData[pathElem].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf(pathElemNotFoundErrTemplate, pathElem)
		}

		return mapTraverse(newMapData, newPath)
	}

	ret, ok := mapData[path[0]]
	if !ok {
		return nil, fmt.Errorf(pathElemNotFoundErrTemplate, path[0])
	}

	return ret, nil
}

func metricTypes(namespace plugin.Namespace, metricsList []Metric) ([]plugin.Metric, error) {
	metrics := []plugin.Metric{}

	for _, m := range metricsList {
		ns := namespace

		for i := range m.Namespace {
			ns = append(ns, m.Namespace[i])
		}

		metrics = append(
			metrics,
			plugin.Metric{
				Namespace: ns,
				Version:   1,
				Tags: map[string]string{
					"Type": m.Type,
					"Path": m.Path,
					"Unit": m.Unit,
				},
			},
		)
	}

	return metrics, nil
}
