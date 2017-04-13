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

// APMMetrics is a map containing the available APM metrics and their properties (only used in this collector).
var APMMetrics = []map[string]string{
	{
		"namespace": "show/health/status",
		"type":      "application",
		"path":      "HealthStatus",
		"unit":      "string",
	},
	{
		"namespace": "show/reporting",
		"type":      "application",
		"path":      "Reporting",
		"unit":      "bool",
	},
	{
		"namespace": "show/summary/application/response_time",
		"type":      "application",
		"path":      "ApplicationSummary/ResponseTime",
		"unit":      "float",
	},
	{
		"namespace": "show/summary/application/throughput",
		"type":      "application",
		"path":      "ApplicationSummary/Throughput",
		"unit":      "float",
	},
	{
		"namespace": "show/summary/application/error_rate",
		"type":      "application",
		"path":      "ApplicationSummary/ErrorRate",
		"unit":      "float",
	},
	{
		"namespace": "show/summary/application/apdex_target",
		"type":      "application",
		"path":      "ApplicationSummary/ApdexTarget",
		"unit":      "float",
	},
	{
		"namespace": "show/summary/application/apdex_score",
		"type":      "application",
		"path":      "ApplicationSummary/ApdexScore",
		"unit":      "float",
	},
	{
		"namespace": "show/summary/application/host_count",
		"type":      "application",
		"path":      "ApplicationSummary/HostCount",
		"unit":      "int",
	},
	{
		"namespace": "show/summary/application/instance_count",
		"type":      "application",
		"path":      "ApplicationSummary/InstanceCount",
		"unit":      "int",
	},
	{
		"namespace": "show/summary/user/response_time",
		"type":      "application",
		"path":      "EndUserSummary/ResponseTime",
		"unit":      "float",
	},
	{
		"namespace": "show/summary/user/throughput",
		"type":      "application",
		"path":      "EndUserSummary/Throughput",
		"unit":      "float",
	},
	{
		"namespace": "show/summary/user/apdex_target",
		"type":      "application",
		"path":      "EndUserSummary/ApdexTarget",
		"unit":      "float",
	},
	{
		"namespace": "show/summary/user/apdex_score",
		"type":      "application",
		"path":      "EndUserSummary/ApdexScore",
		"unit":      "float",
	},
}

// APMClient is the interface every AMP client needs to implement.
type APMClient interface {
	GetApplication(int) (*nr.Application, error)
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
		ns = ns.AddStaticElements(strings.Split(m["namespace"], "/")...)

		metrics = append(
			metrics,
			plugin.Metric{
				Namespace: ns,
				Version:   1,
				Tags:      m,
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

	apps := map[int]*nr.Application{}
	for i, m := range metrics {
		if m.Namespace.Element(2).Value != "apm" {
			continue
		}

		appID := m.Namespace.Element(3)

		appIDInt, err := strconv.Atoi(appID.Value)
		if err != nil {
			return collectedMetrics, err
		}

		if _, ok := apps[appIDInt]; !ok {
			// Application info missing, fetching...
			app, err := a.APMClient.GetApplication(appIDInt)
			if err != nil {
				return collectedMetrics, err
			}

			apps[appIDInt] = app
		}

		mapData := map[string]interface{}{}
		switch m.Tags["type"] {
		case "application":
			appInfo := apps[appIDInt]
			mapData = structs.Map(appInfo)

			break
		default:
			return collectedMetrics, fmt.Errorf("Unsupported New Relic data type")
		}

		// Create a new metric based on the "old" one.
		newMetric := metrics[i]

		metricData, err := mapTraverse(mapData, strings.Split(m.Tags["path"], "/"))
		if err != nil {
			return collectedMetrics, err
		}

		newMetric.Data = metricData
		newMetric.Unit = m.Tags["unit"]
		newMetric.Tags = map[string]string{}
		newMetric.Timestamp = time.Now().UTC()

		collectedMetrics = append(
			collectedMetrics,
			newMetric,
		)
	}

	return collectedMetrics, nil
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
