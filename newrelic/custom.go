// This is in fact the metrics object (metric data). The custom name is just to avoid collission with all the "metric"
// references all around in the code.

package newrelic

import (
	"fmt"
	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
	nr "github.com/yfronto/newrelic"
	"strconv"
	"time"
)

// CustomMetrics defines the available metric data metrics.
var CustomMetrics = []Metric{
	{
		Namespace: plugin.Namespace{
			plugin.NewNamespaceElement("application"),
			plugin.NamespaceElement{
				Name:        "app_id",
				Description: "The application id",
				Value:       "*",
			},
			plugin.NamespaceElement{
				Name:        "minutes",
				Description: "Number of minutes to construct a relative timeframe from (now - minutes).",
				Value:       "*",
			},
			plugin.NamespaceElement{
				Name:        "metric_name",
				Description: "Metric name",
				Value:       "*",
			},
			plugin.NamespaceElement{
				Name:        "value_name",
				Description: "Value name",
				Value:       "*",
			},
			plugin.NewNamespaceElement("value"),
		},
		Type: "application",
		Unit: "float",
	},
	{
		Namespace: plugin.Namespace{
			plugin.NewNamespaceElement("component"),
			plugin.NamespaceElement{
				Name:        "component_id",
				Description: "The component id",
				Value:       "*",
			},
			plugin.NamespaceElement{
				Name:        "minutes",
				Description: "Number of minutes to construct a relative timeframe from (now - minutes).",
				Value:       "*",
			},
			plugin.NamespaceElement{
				Name:        "metric_name",
				Description: "Metric name",
				Value:       "*",
			},
			plugin.NamespaceElement{
				Name:        "value_name",
				Description: "Value name",
				Value:       "*",
			},
			plugin.NewNamespaceElement("value"),
		},
		Type: "component",
		Unit: "float",
	},
}

// CustomClient defines the custom metrics (all metric data metrics) client.
type CustomClient interface {
	GetApplicationMetricData(int, []string, *nr.MetricDataOptions) (*nr.MetricDataResponse, error)
	GetComponentMetricData(int, []string, *nr.MetricDataOptions) (*nr.MetricDataResponse, error)
}

// CustomClientImpl is a real implementation of an CustomClient.
type CustomClientImpl struct {
	APIKey string
}

// GetApplicationMetricData fetches application specific metric data.
func (cc *CustomClientImpl) GetApplicationMetricData(appID int, names []string, options *nr.MetricDataOptions) (*nr.MetricDataResponse, error) {
	c := nr.NewClient(cc.APIKey)

	return c.GetApplicationMetricData(appID, names, options)
}

// GetComponentMetricData fetches component specific metric data.
func (cc *CustomClientImpl) GetComponentMetricData(componentID int, names []string, options *nr.MetricDataOptions) (*nr.MetricDataResponse, error) {
	c := nr.NewClient(cc.APIKey)

	return c.GetComponentMetricData(componentID, names, options)
}

// Custom represents the custom metric data metrics available from New Relic.
type Custom struct {
	CustomClient CustomClient
}

// NewCustom creates and returns a new Custom object with a configured CustomClient.
func NewCustom(apiKey string) Service {
	return &Custom{
		CustomClient: &CustomClientImpl{
			APIKey: apiKey,
		},
	}
}

// GetMetricTypes returns the available metric types.
func (c *Custom) GetMetricTypes(_ plugin.Config) ([]plugin.Metric, error) {
	ns := plugin.NewNamespace("inteleon", "newrelic", "metric")

	return metricTypes(ns, CustomMetrics)
}

// CollectMetrics fetches the requested metric data metrics and returns them.
func (c *Custom) CollectMetrics(metrics []plugin.Metric) ([]plugin.Metric, error) {
	collectedMetrics := []plugin.Metric{}

	metricResponses := map[string]map[int]map[string]*nr.MetricDataResponse{}
	for i, m := range metrics {
		if m.Namespace.Element(2).Value != "metric" {
			continue
		}

		metricType := m.Tags["Type"]
		id := m.Namespace.Element(4)

		idInt, err := strconv.Atoi(id.Value)
		if err != nil {
			return collectedMetrics, err
		}

		relativeMin := m.Namespace.Element(5).Value
		metricStringID := m.Namespace.Element(6).Value
		metricDataOptions := &nr.MetricDataOptions{
			Summarize: true,
		}

		if relativeMin != "*" {
			relativeMinInt, err := strconv.Atoi(relativeMin)
			if err != nil {
				return collectedMetrics, err
			}

			metricDataOptions.From = time.Now().UTC().Add(-(time.Duration(relativeMinInt) * time.Minute))
			metricDataOptions.To = time.Now().UTC()
		}

		if _, ok := metricResponses[metricType]; !ok {
			metricResponses[metricType] = map[int]map[string]*nr.MetricDataResponse{}
		}

		if _, ok := metricResponses[metricType][idInt][metricStringID]; !ok {
			// Metrics missing, fetching...
			var fetchMetricData *nr.MetricDataResponse
			var err error

			switch metricType {
			case "application":
				fetchMetricData, err = c.CustomClient.GetApplicationMetricData(
					idInt,
					[]string{metricStringID},
					metricDataOptions,
				)

				break
			case "component":
				fetchMetricData, err = c.CustomClient.GetComponentMetricData(
					idInt,
					[]string{metricStringID},
					metricDataOptions,
				)

				break
			}

			if err != nil {
				return collectedMetrics, err
			}

			if _, ok := metricResponses[metricType][idInt]; !ok {
				metricResponses[metricType][idInt] = map[string]*nr.MetricDataResponse{}
			}

			metricResponses[metricType][idInt][metricStringID] = fetchMetricData
		}

		metricData := metricResponses[metricType][idInt][metricStringID]

		numberOfMetricsFound := len(metricData.Metrics)
		if numberOfMetricsFound != 1 {
			// Metric not found, skip reporting it and continue execution.
			continue
		}

		firstMetric := metricData.Metrics[0]
		if firstMetric.Name != metricStringID {
			return collectedMetrics, fmt.Errorf(
				"Metric name mismatch! Requested metric name: %s. Metric name in the received payload: %s.",
				metricStringID,
				firstMetric.Name,
			)
		}

		metricValues := firstMetric.Timeslices[0].Values
		castValues := map[string]interface{}{}
		for ci := range metricValues {
			castValues[ci] = metricValues[ci]
		}

		populatedMetric, err := populateMetric(metrics[i], castValues)
		if err != nil {
			return collectedMetrics, err
		}

		collectedMetrics = append(collectedMetrics, populatedMetric)
	}

	return collectedMetrics, nil
}
