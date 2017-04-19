package newrelic

import (
	"fmt"
	"github.com/fatih/structs"
	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
	nr "github.com/yfronto/newrelic"
	"strconv"
	"strings"
	"time"
)

// APMMetrics is a list containing the available APM metrics and their properties.
var APMMetrics = []Metric{
	{
		Namespace: plugin.NewNamespace("show", "health", "status"),
		Type:      "application",
		Path:      "HealthStatus",
		Unit:      "string",
	},
	{
		Namespace: plugin.NewNamespace("show", "reporting"),
		Type:      "application",
		Path:      "Reporting",
		Unit:      "bool",
	},
	{
		Namespace: plugin.NewNamespace("show", "summary", "application", "response_time"),
		Type:      "application",
		Path:      "ApplicationSummary/ResponseTime",
		Unit:      "float",
	},
	{
		Namespace: plugin.NewNamespace("show", "summary", "application", "throughput"),
		Type:      "application",
		Path:      "ApplicationSummary/Throughput",
		Unit:      "float",
	},
	{
		Namespace: plugin.NewNamespace("show", "summary", "application", "error_rate"),
		Type:      "application",
		Path:      "ApplicationSummary/ErrorRate",
		Unit:      "float",
	},
	{
		Namespace: plugin.NewNamespace("show", "summary", "application", "apdex_target"),
		Type:      "application",
		Path:      "ApplicationSummary/ApdexTarget",
		Unit:      "float",
	},
	{
		Namespace: plugin.NewNamespace("show", "summary", "application", "apdex_score"),
		Type:      "application",
		Path:      "ApplicationSummary/ApdexScore",
		Unit:      "float",
	},
	{
		Namespace: plugin.NewNamespace("show", "summary", "application", "host_count"),
		Type:      "application",
		Path:      "ApplicationSummary/HostCount",
		Unit:      "int",
	},
	{
		Namespace: plugin.NewNamespace("show", "summary", "application", "instance_count"),
		Type:      "application",
		Path:      "ApplicationSummary/InstanceCount",
		Unit:      "int",
	},
	{
		Namespace: plugin.NewNamespace("show", "summary", "user", "response_time"),
		Type:      "application",
		Path:      "EndUserSummary/ResponseTime",
		Unit:      "float",
	},
	{
		Namespace: plugin.NewNamespace("show", "summary", "user", "throughput"),
		Type:      "application",
		Path:      "EndUserSummary/Throughput",
		Unit:      "float",
	},
	{
		Namespace: plugin.NewNamespace("show", "summary", "user", "apdex_target"),
		Type:      "application",
		Path:      "EndUserSummary/ApdexTarget",
		Unit:      "float",
	},
	{
		Namespace: plugin.NewNamespace("show", "summary", "user", "apdex_score"),
		Type:      "application",
		Path:      "EndUserSummary/ApdexScore",
		Unit:      "float",
	},
	{
		Namespace: plugin.Namespace{
			plugin.NewNamespaceElement("metric"),
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
		},
		Type: "metric",
		Unit: "float",
	},
}

// APMClient is the interface every AMP client needs to implement.
type APMClient interface {
	GetApplication(int) (*nr.Application, error)
	GetApplicationMetricData(int, []string, *nr.MetricDataOptions) (*nr.MetricDataResponse, error)
}

// APMClientImpl is a real implementation of an APMClient.
type APMClientImpl struct {
	APIKey string
}

// GetApplication fetches application information from New Relic (APM).
func (a *APMClientImpl) GetApplication(appID int) (*nr.Application, error) {
	c := nr.NewClient(a.APIKey)

	return c.GetApplication(appID)
}

// GetApplicationMetricData fetches application specific metric data (custom or not).
func (a *APMClientImpl) GetApplicationMetricData(appID int, names []string, options *nr.MetricDataOptions) (*nr.MetricDataResponse, error) {
	c := nr.NewClient(a.APIKey)

	return c.GetApplicationMetricData(appID, names, options)
}

// APM represents the APM service part of New Relic.
type APM struct {
	APMClient APMClient
}

// NewAPM creates and returns a new APM object with a configured APMClient.
func NewAPM(apiKey string) *APM {
	return &APM{
		APMClient: &APMClientImpl{
			APIKey: apiKey,
		},
	}
}

// GetMetricTypes returns the available APM metric types.
func (a *APM) GetMetricTypes(_ plugin.Config) ([]plugin.Metric, error) {
	metrics := []plugin.Metric{}

	for _, m := range APMMetrics {
		ns := plugin.NewNamespace("inteleon", "newrelic", "apm")
		ns = ns.AddDynamicElement("app_id", "Application id")

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

// CollectMetrics fetches the requested APM metrics and returns them.
func (a *APM) CollectMetrics(metrics []plugin.Metric) ([]plugin.Metric, error) {
	collectedMetrics := []plugin.Metric{}

	if len(metrics) == 0 {
		return collectedMetrics, fmt.Errorf("List of metrics is empty")
	}

	apps := []plugin.Metric{}
	appAdditionalMetrics := []plugin.Metric{}
	for i, m := range metrics {
		if m.Namespace.Element(2).Value != "apm" {
			continue
		}

		switch m.Tags["Type"] {
		case "application":
			apps = append(apps, metrics[i])

			break
		case "metric":
			appAdditionalMetrics = append(appAdditionalMetrics, metrics[i])

			break
		}
	}

	if len(apps) > 0 {
		appsMetrics, err := a.collectApplications(apps)
		if err != nil {
			return collectedMetrics, err
		}

		for i := range appsMetrics {
			collectedMetrics = append(collectedMetrics, appsMetrics[i])
		}
	}

	if len(appAdditionalMetrics) > 0 {
		appAdditionalMetricsData, err := a.collectApplicationAdditionalMetrics(appAdditionalMetrics)
		if err != nil {
			return collectedMetrics, err
		}

		for i := range appAdditionalMetricsData {
			collectedMetrics = append(collectedMetrics, appAdditionalMetricsData[i])
		}
	}

	return collectedMetrics, nil
}

func (a *APM) populateMetric(metric plugin.Metric, mapData map[string]interface{}) (plugin.Metric, error) {
	// Create a new metric based on the "old" one.
	newMetric := metric

	mPath := []string{}
	if metric.Tags["Path"] == "" {
		mPath = []string{metric.Namespace.Element(len(metric.Namespace) - 1).Value}
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

func (a *APM) collectApplications(metrics []plugin.Metric) ([]plugin.Metric, error) {
	appsMetrics := []plugin.Metric{}

	apps := map[int]*nr.Application{}
	for i, m := range metrics {
		appID := m.Namespace.Element(3)

		appIDInt, err := strconv.Atoi(appID.Value)
		if err != nil {
			return appsMetrics, err
		}

		if _, ok := apps[appIDInt]; !ok {
			// Application info missing, fetching...
			app, err := a.APMClient.GetApplication(appIDInt)
			if err != nil {
				return appsMetrics, err
			}

			apps[appIDInt] = app
		}

		// Convert the app data to a struct so it's more easily traversable and more universal before passing it to the populateMetric function.
		appMetric, err := a.populateMetric(metrics[i], structs.Map(apps[appIDInt]))
		if err != nil {
			return appsMetrics, err
		}

		appsMetrics = append(appsMetrics, appMetric)
	}

	return appsMetrics, nil
}

func (a *APM) collectApplicationAdditionalMetrics(metrics []plugin.Metric) ([]plugin.Metric, error) {
	appAdditionalMetrics := []plugin.Metric{}

	for i, m := range metrics {
		appID := m.Namespace.Element(3)

		appIDInt, err := strconv.Atoi(appID.Value)
		if err != nil {
			return appAdditionalMetrics, err
		}

		metricStringID := m.Namespace.Element(5)
		appAdditionalMetricData, err := a.APMClient.GetApplicationMetricData(appIDInt, []string{metricStringID.Value}, nil)
		if err != nil {
			return appAdditionalMetrics, err
		}

		numberOfMetricsFound := len(appAdditionalMetricData.Metrics)
		if numberOfMetricsFound != 1 {
			return appAdditionalMetrics, fmt.Errorf(
				"Wrong number of returned metrics when fetching data for %s. Exepected number is 1, got %d.",
				metricStringID,
				numberOfMetricsFound,
			)
		}

		firstMetric := appAdditionalMetricData.Metrics[0]
		if firstMetric.Name != metricStringID.Value {
			return appAdditionalMetrics, fmt.Errorf(
				"Metric name mismatch! Requested metric name: %s. Metric name in the received payload: %s.",
				metricStringID.Value,
				firstMetric.Name,
			)
		}

		metricValues := firstMetric.Timeslices[0].Values
		castValues := map[string]interface{}{}
		for ci := range metricValues {
			castValues[ci] = metricValues[ci]
		}

		appMetric, err := a.populateMetric(metrics[i], castValues)
		if err != nil {
			return appAdditionalMetrics, err
		}

		appAdditionalMetrics = append(appAdditionalMetrics, appMetric)
	}

	return appAdditionalMetrics, nil
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
